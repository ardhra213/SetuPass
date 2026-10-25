package com.setupass.identity

import android.content.Context
import android.net.Uri
import android.os.Handler
import android.os.Looper
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import org.json.JSONArray
import org.json.JSONObject
import java.io.File
import java.net.HttpURLConnection
import java.net.URL
import java.net.URLEncoder
import java.security.MessageDigest
import java.text.SimpleDateFormat
import java.util.*
import java.util.concurrent.Executors

/** Local hackathon login only. Production enrollment must be server-authorized. */
internal object SetuLogin {
    const val ISSUER = "http://localhost:8088"
    const val REDIRECT = "com.setupass.identity://oauth/callback"
    private const val CLIENT_KEY = "setupass.client.native.v2"
    private const val DPOP_KEY = "setupass.dpop.native.v2"
    private val worker = Executors.newSingleThreadExecutor()
    private val main = Handler(Looper.getMainLooper())
    private lateinit var app: Context
    private var accessToken: String? = null
    private var sessionExpiry = 0L
    private val prefs get() = app.getSharedPreferences("setupass.registration", Context.MODE_PRIVATE)
    private val pendingFile get() = File(app.noBackupFilesDir, "login-pending.enc")
    var status by mutableStateOf("Sign in to start your Setu ID session."); private set
    var busy by mutableStateOf(false); private set
    var signedIn by mutableStateOf(false); private set
    var subject by mutableStateOf(""); private set
    fun initialize(context: Context) {
        app = context.applicationContext
        worker.execute {
            try {
                register()
                val saved = try { NativeSession.get(app) } catch (_: Exception) { null }
                main.post {
                    if (saved != null) {
                        subject = saved.getString("sub"); signedIn = true; sessionExpiry = saved.getLong("expires")
                        status = "Native SSO session available. Open Setu Benefits to continue."
                    } else status = "Ready. After opening both apps, run 03_PAIR_APPS.cmd once, then sign in."
                }
            } catch (e: Exception) { fail(e) }
        }
    }
    private fun now() = System.currentTimeMillis() / 1000
    private fun note(text: String) { main.post { status = text } }
    private fun encode(text: String) = URLEncoder.encode(text, "UTF-8")
    private fun clientId() = "setupass-native-parent-" + AuthCrypto.thumbprint(AuthCrypto.publicKey(CLIENT_KEY)).take(20)
    private fun assertion(): String {
        val id = clientId(); val jwk = AuthCrypto.publicKey(CLIENT_KEY)
        return AuthCrypto.sign(CLIENT_KEY,
            JSONObject().put("typ", "JWT").put("alg", "ES256").put("kid", AuthCrypto.thumbprint(jwk)),
            JSONObject().put("iss", id).put("sub", id).put("aud", ISSUER)
                .put("iat", now()).put("exp", now() + 60).put("jti", UUID.randomUUID().toString()))
    }
    internal fun proof(url: String, method: String, token: String? = null, nonce: String? = null): String {
        val claims = JSONObject().put("jti", UUID.randomUUID().toString()).put("htm", method)
            .put("htu", url).put("iat", now())
        if (token != null) claims.put("ath", AuthCrypto.hash(token))
        if (nonce != null) claims.put("nonce", nonce)
        return AuthCrypto.sign(DPOP_KEY,
            JSONObject().put("typ", "dpop+jwt").put("alg", "ES256").put("jwk", AuthCrypto.publicKey(DPOP_KEY)), claims)
    }
    private fun request(path: String, body: String? = null, form: Boolean = false, dpop: Boolean = false): JSONObject {
        var serverNonce: String? = null
        repeat(2) { attempt ->
            val url = ISSUER + path
            val conn = URL(url).openConnection() as HttpURLConnection
            try {
                conn.connectTimeout = 15000; conn.readTimeout = 20000; conn.instanceFollowRedirects = false
                conn.requestMethod = if (body == null) "GET" else "POST"
                conn.setRequestProperty("Accept", "application/json")
                if (dpop) conn.setRequestProperty("DPoP", proof(url, conn.requestMethod, nonce = serverNonce))
                if (body != null) {
                    conn.doOutput = true
                    conn.setRequestProperty("Content-Type", if (form) "application/x-www-form-urlencoded" else "application/json")
                    conn.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
                }
                val code = conn.responseCode
                val stream = if (code in 200..299) conn.inputStream else conn.errorStream
                val raw = stream?.bufferedReader()?.use { it.readText() } ?: ""
                val json = try { JSONObject(raw) } catch (_: Exception) { JSONObject() }
                val nonce = conn.getHeaderField("DPoP-Nonce")
                if (attempt == 0 && dpop && nonce != null && json.optString("error") == "use_dpop_nonce") {
                    serverNonce = nonce
                } else {
                    val errors = json.optJSONArray("errors")
                    val error = if (errors != null && errors.length() > 0) errors.getJSONObject(0).optString("errorCode") else json.optString("error")
                    require(code in 200..299 && error.isEmpty()) {
                        // No JWT, authorization code, private key, or response body is shown in diagnostics.
                        "$path: HTTP $code ${error.take(100).replace(Regex("[^A-Za-z0-9_ .-]"), "")}".trim()
                    }
                    require(json.length() > 0) { "$path returned no JSON" }
                    return json
                }
            } finally { conn.disconnect() }
        }
        error("Server rejected the DPoP nonce retry")
    }
    private fun form(path: String, values: Map<String, String>): JSONObject = request(path,
        values.entries.joinToString("&") { encode(it.key) + "=" + encode(it.value) }, form = true, dpop = true)
    private fun register() {
        val id = clientId()
        if (prefs.getString("registered", null) == id) return
        note("Setting up this phone for Setu ID…")
        val jwk = AuthCrypto.publicKey(CLIENT_KEY).put("kid", AuthCrypto.thumbprint(AuthCrypto.publicKey(CLIENT_KEY))).put("use", "sig")
        val registration = JSONObject().put("clientId", id).put("clientName", "Setu ID")
            .put("clientNameLangMap", JSONObject().put("eng", "Setu ID"))
            .put("relyingPartyId", "decl-ou-1").put("logoUri", "http://localhost:3000/logo.png")
            .put("redirectUris", JSONArray().put(REDIRECT))
            .put("userClaims", JSONArray().put("name").put("email").put("picture"))
            .put("authContextRefs", JSONArray().put("mosip:idp:acr:generated-code"))
            .put("publicKey", jwk).put("grantTypes", JSONArray().put("authorization_code"))
            .put("clientAuthMethods", JSONArray().put("private_key_jwt"))
            .put("additionalConfig", JSONObject().put("require_pushed_authorization_requests", true)
                .put("dpop_bound_access_tokens", true).put("require_pkce", true)
                .put("allowed_authorization_scopes", JSONArray().put("device_sso")))
        val date = SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'", Locale.US).apply { timeZone = TimeZone.getTimeZone("UTC") }.format(Date())
        val result = request("/client-mgmt/client", JSONObject().put("requestTime", date).put("request", registration).toString())
        val response = result.getJSONObject("response")
        require(response.getString("clientId") == id && response.getString("status") == "ACTIVE") { "Client registration did not become ACTIVE" }
        check(prefs.edit().putString("registered", id).commit()) { "Could not save registration" }
    }
    fun start(openBrowser: (Uri) -> Unit) {
        if (busy) return
        busy = true; signedIn = false; subject = ""; accessToken = null; status = "Connecting to eSignet…"
        worker.execute {
            try {
                pendingFile.delete()
                NativeSession.clear(app)
                register()
                val state = AuthCrypto.random(); val verifier = AuthCrypto.random(); val nonce = AuthCrypto.random()
                note("Preparing secure login…")
                val result = form("/oauth2/par", mapOf(
                    "client_id" to clientId(), "redirect_uri" to REDIRECT, "response_type" to "code",
                    "scope" to "openid profile email device_sso", "state" to state, "nonce" to nonce,
                    "code_challenge" to AuthCrypto.hash(verifier), "code_challenge_method" to "S256",
                    "client_assertion_type" to "urn:ietf:params:oauth:client-assertion-type:jwt-bearer",
                    "client_assertion" to assertion(), "dpop_jkt" to AuthCrypto.thumbprint(AuthCrypto.publicKey(DPOP_KEY))))
                val pending = JSONObject().put("state", state).put("verifier", verifier).put("nonce", nonce).put("created", now())
                pendingFile.writeBytes(AuthCrypto.encrypt(pending.toString()))
                val uri = Uri.parse("$ISSUER/oauth2/authorize").buildUpon()
                    .appendQueryParameter("client_id", clientId()).appendQueryParameter("request_uri", result.getString("request_uri")).build()
                main.post {
                    busy = false; status = "Finish OTP login and consent in your browser."
                    try { openBrowser(uri) } catch (_: Exception) { status = "No browser opened. Install or enable Chrome, then try again." }
                }
            } catch (e: Exception) { fail(e) }
        }
    }
    fun receive(uri: Uri) {
        if (uri.scheme != "com.setupass.identity" || uri.host != "oauth" || uri.path != "/callback") return
        if (busy) return
        busy = true; status = "Finishing sign-in…"
        worker.execute {
            try {
                require(pendingFile.exists()) { "No pending login. Tap Sign in to start again." }
                val pending = JSONObject(AuthCrypto.decrypt(pendingFile.readBytes()))
                require(now() - pending.getLong("created") in 0..3600) { "Login expired. Start again." }
                require(uri.getQueryParameters("state").size == 1) { "Missing or duplicated login state" }
                val state = uri.getQueryParameter("state") ?: ""
                require(MessageDigest.isEqual(state.toByteArray(), pending.getString("state").toByteArray())) { "Login state did not match. Start again." }
                val returnIssuer = uri.getQueryParameter("iss")
                require(returnIssuer == null || returnIssuer == ISSUER) { "Unexpected login issuer" }
                pendingFile.delete() // Consume once; failed exchange requires a new login.
                require(uri.getQueryParameter("error") == null) { "Browser login was declined or expired. Start again." }
                require(uri.getQueryParameters("code").size == 1) { "Missing authorization code. Start again." }
                val code = uri.getQueryParameter("code") ?: ""
                require(code.isNotBlank()) { "Missing authorization code" }
                val tokens = form("/oauth2/token", mapOf(
                    "grant_type" to "authorization_code", "code" to code, "redirect_uri" to REDIRECT,
                    "client_id" to clientId(), "code_verifier" to pending.getString("verifier"),
                    "client_assertion_type" to "urn:ietf:params:oauth:client-assertion-type:jwt-bearer", "client_assertion" to assertion()))
                require(tokens.getString("token_type").equals("DPoP", true)) { "Expected a DPoP-bound token" }
                val claims = AuthCrypto.verifyIdToken(tokens.getString("id_token"), request("/oauth2/jwks"))
                validateClaims(claims, pending.getString("nonce"))
                val ttl = tokens.getLong("expires_in")
                require(ttl in 1..86400) { "Unexpected token lifetime" }
                val token = tokens.getString("access_token"); require(token.isNotBlank())
                val secret = tokens.getString("device_secret")
                require(claims.getString("ds_hash") == AuthCrypto.hash(secret)) { "Device secret binding failed" }
                require(claims.getString("sid").isNotBlank()) { "Session ID missing" }
                NativeSession.save(app, tokens.getString("id_token"), secret, claims.getLong("exp"), claims.getString("sub"))
                val expiry = claims.getLong("exp")
                main.post {
                    accessToken = token; sessionExpiry = expiry; subject = claims.getString("sub")
                    signedIn = true; busy = false; status = "Native SSO is ready. Open Setu Benefits and continue without another OTP. Session lasts up to 15 minutes."
                    main.postDelayed({
                        if (sessionExpiry == expiry && now() >= expiry) {
                            signedIn = false; accessToken = null; subject = ""
                            status = "This short demo session expired. Sign in again to continue."
                        }
                    }, ((expiry - now()).coerceAtLeast(0) * 1000) + 1000)
                }
            } catch (e: Exception) { fail(e) }
        }
    }
    private fun validateClaims(c: JSONObject, nonce: String) {
        require(c.getString("iss") == ISSUER) { "ID token issuer mismatch" }
        val aud = c.get("aud")
        val audiences = if (aud is JSONArray) (0 until aud.length()).map { aud.getString(it) } else listOf(aud.toString())
        require(clientId() in audiences) { "ID token audience mismatch" }
        if (audiences.size > 1 || c.has("azp")) require(c.optString("azp") == clientId()) { "ID token authorized party mismatch" }
        require(c.getLong("exp") > now() && c.getLong("iat") <= now() + 30) { "ID token expired or phone clock is incorrect" }
        if (c.has("nbf")) require(c.getLong("nbf") <= now() + 30) { "ID token not yet valid" }
        require(c.optString("nonce") == nonce) { "ID token nonce mismatch" }
        require(c.getString("sub").isNotBlank()) { "ID token subject missing" }
    }
    private fun fail(e: Exception) {
        main.post {
            busy = false; signedIn = false; accessToken = null; subject = ""
            status = when (e) {
                is java.net.ConnectException, is java.net.SocketTimeoutException -> "Cannot reach eSignet. Check Docker and the two USB port commands, then try again."
                is IllegalArgumentException, is IllegalStateException -> e.message ?: "Login validation failed. Start again."
                else -> "Login stopped (${e.javaClass.simpleName}). Share this message so we can fix it."
            }
        }
    }
    fun clearSession() {
        if (busy) return
        signedIn = false; accessToken = null; subject = ""; sessionExpiry = 0
        worker.execute { pendingFile.delete(); NativeSession.clear(app) }
        status = "Session cleared on this app. Your browser session may still be active."
    }
}
