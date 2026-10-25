package com.setupass.benefits

import android.content.Context
import android.net.Uri
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import org.json.JSONArray
import org.json.JSONObject
import java.net.HttpURLConnection
import java.net.URL
import java.net.URLEncoder
import java.text.SimpleDateFormat
import java.util.*
import java.util.concurrent.Executors

internal object BenefitsLogin {
    private const val ISSUER = "http://localhost:8088"
    private const val KEY = "setupass.benefits.client.v1"
    private val worker = Executors.newSingleThreadExecutor()
    private val main = Handler(Looper.getMainLooper())
    private lateinit var app: Context
    var status by mutableStateOf("Preparing this app…"); private set
    var busy by mutableStateOf(false); private set
    var subject by mutableStateOf(""); private set
    var signedIn by mutableStateOf(false); private set
    private var sessionGeneration = 0L
    private fun now() = System.currentTimeMillis() / 1000
    private fun id() = "setupass-native-child-" + AuthCrypto.thumbprint(AuthCrypto.publicKey(KEY)).take(20)
    fun initialize(c: Context) {
        app = c.applicationContext
        worker.execute { try { register(); main.post { status = "Ready. Sign in to Setu ID, then continue here." } } catch (e: Exception) { fail(e) } }
    }
    private fun assertion(): String = AuthCrypto.sign(KEY,
        JSONObject().put("typ", "JWT").put("alg", "ES256").put("kid", AuthCrypto.thumbprint(AuthCrypto.publicKey(KEY))),
        JSONObject().put("iss", id()).put("sub", id()).put("aud", ISSUER).put("iat", now()).put("exp", now()+60).put("jti", UUID.randomUUID().toString()))
    private fun broker(method: String, nonce: String? = null): Bundle {
        val b = app.contentResolver.call(Uri.parse("content://com.setupass.identity.sso"), method, null,
            Bundle().apply { if (nonce != null) putString("nonce", nonce) }) ?: error("Install Setu ID first.")
        check(!b.containsKey("error")) { b.getString("error") ?: "Parent session unavailable" }
        return b
    }
    private fun request(path: String, body: String? = null, form: Boolean = false): JSONObject {
        var nonce: String? = null
        repeat(2) { attempt ->
            val c = URL(ISSUER+path).openConnection() as HttpURLConnection
            try {
                c.connectTimeout=15000;c.readTimeout=20000;c.instanceFollowRedirects=false
                c.requestMethod=if(body==null) "GET" else "POST"
                if(form) c.setRequestProperty("DPoP", requireNotNull(broker("proof",nonce).getString("dpop")))
                if(body!=null) {
                    c.doOutput=true;c.setRequestProperty("Content-Type",if(form) "application/x-www-form-urlencoded" else "application/json")
                    c.outputStream.use { it.write(body.toByteArray(Charsets.UTF_8)) }
                }
                val code=c.responseCode
                val raw=(if(code in 200..299)c.inputStream else c.errorStream)?.bufferedReader()?.use{it.readText()}?:""
                val json=try{JSONObject(raw)}catch(_:Exception){JSONObject()}
                val next=c.getHeaderField("DPoP-Nonce")
                if(attempt==0 && form && next!=null && json.optString("error")=="use_dpop_nonce") { nonce=next }
                else {
                    val errors=json.optJSONArray("errors")
                    val error=if(errors!=null && errors.length()>0)errors.getJSONObject(0).optString("errorCode") else json.optString("error")
                    check(code in 200..299 && error.isEmpty()) { "HTTP $code ${error.take(80).replace(Regex("[^A-Za-z0-9_-]"), "")}. If pairing is pending, run 03_PAIR_APPS.cmd." }
                    return json
                }
            } finally { c.disconnect() }
        }
        error("DPoP nonce retry failed")
    }
    private fun register() {
        val prefs=app.getSharedPreferences("registration",Context.MODE_PRIVATE)
        if(prefs.getString("id",null)==id())return
        val jwk=AuthCrypto.publicKey(KEY).put("kid",AuthCrypto.thumbprint(AuthCrypto.publicKey(KEY))).put("use","sig")
        val registration=JSONObject().put("clientId",id()).put("clientName","Setu Benefits")
            .put("clientNameLangMap",JSONObject().put("eng","Setu Benefits")).put("relyingPartyId","decl-ou-1")
            .put("logoUri","http://localhost:3000/logo.png").put("redirectUris",JSONArray().put("com.setupass.benefits://oauth/callback"))
            .put("userClaims",JSONArray().put("name").put("email"))
            .put("authContextRefs",JSONArray().put("mosip:idp:acr:generated-code"))
            .put("publicKey",jwk).put("grantTypes",JSONArray().put("urn:ietf:params:oauth:grant-type:token-exchange"))
            .put("clientAuthMethods",JSONArray().put("private_key_jwt"))
            .put("additionalConfig",JSONObject().put("dpop_bound_access_tokens",true).put("require_pkce",false))
        val date=SimpleDateFormat("yyyy-MM-dd'T'HH:mm:ss.SSS'Z'",Locale.US).apply{timeZone=TimeZone.getTimeZone("UTC")}.format(Date())
        val r=request("/client-mgmt/client",JSONObject().put("requestTime",date).put("request",registration).toString()).getJSONObject("response")
        check(r.getString("clientId")==id() && r.getString("status")=="ACTIVE")
        check(prefs.edit().putString("id",id()).commit())
    }
    fun continueSSO() {
        if(busy)return
        busy=true;signedIn=false;subject="";status="Asking Setu ID for your device session…"
        worker.execute {
            try {
                register()
                val shared=broker("credentials")
                val values=mapOf("grant_type" to "urn:ietf:params:oauth:grant-type:token-exchange",
                    "client_id" to id(),"scope" to "openid profile email",
                    "subject_token_type" to "urn:ietf:params:oauth:token-type:id_token",
                    "subject_token" to requireNotNull(shared.getString("id_token")),
                    "actor_token_type" to "urn:openid:params:token-type:device-secret",
                    "actor_token" to requireNotNull(shared.getString("device_secret")),
                    "client_assertion_type" to "urn:ietf:params:oauth:client-assertion-type:jwt-bearer","client_assertion" to assertion())
                val body=values.entries.joinToString("&"){URLEncoder.encode(it.key,"UTF-8")+"="+URLEncoder.encode(it.value,"UTF-8")}
                val tokens=request("/oauth2/token",body,true)
                check(tokens.getString("token_type").equals("DPoP",true)) { "Missing device-bound token" }
                val c=AuthCrypto.verifyIdToken(tokens.getString("id_token"),request("/oauth2/jwks"))
                val aud=c.get("aud");val audiences=if(aud is JSONArray)(0 until aud.length()).map{aud.getString(it)} else listOf(aud.toString())
                check(c.getString("iss")==ISSUER && audiences==listOf(id()) && c.getLong("exp")>now() && c.getLong("iat")<=now()+30) { "Token validation failed" }
                check(!c.has("nbf") || c.getLong("nbf")<=now()+30)
                check(!c.has("azp") || c.getString("azp")==id())
                check(c.getString("sub")==shared.getString("subject") && c.getString("sid").isNotBlank()) { "Identity binding failed" }
                check(!tokens.has("device_secret")) { "Child must not receive a new device secret" }
                main.post {
                    signedIn=true;busy=false;subject=c.getString("sub")
                    val generation=++sessionGeneration
                    status="Signed in through Setu ID. No browser, password or second OTP."
                    main.postDelayed({if(sessionGeneration==generation){signedIn=false;subject="";status="App session expired. Continue with Setu ID again."}},((c.getLong("exp")-now()).coerceAtLeast(0)*1000)+1000)
                }
            } catch(e:Exception){fail(e)}
        }
    }
    private fun fail(e:Exception) { main.post { busy=false;signedIn=false;subject="";status=if(e is IllegalStateException)e.message?:"Sign in again" else "Connection or session unavailable. Check Docker, USB, pairing, and Setu ID login. (${e.javaClass.simpleName})" } }
    fun clear(){sessionGeneration++;signedIn=false;subject="";status="Benefits session cleared. Your Setu ID session is separate."}
}
