package com.setupass.identity

import android.security.keystore.KeyGenParameterSpec
import android.security.keystore.KeyProperties
import android.util.Base64
import org.json.JSONObject
import java.math.BigInteger
import java.security.*
import java.security.interfaces.ECPublicKey
import java.security.spec.*
import javax.crypto.Cipher
import javax.crypto.KeyGenerator
import javax.crypto.SecretKey
import javax.crypto.spec.GCMParameterSpec

/** Device keys stay in Android Keystore. No Postman key is imported. */
internal object AuthCrypto {
    private fun store() = KeyStore.getInstance("AndroidKeyStore").apply { load(null) }
    fun b64(bytes: ByteArray): String = Base64.encodeToString(bytes, Base64.URL_SAFE or Base64.NO_PADDING or Base64.NO_WRAP)
    fun decode(s: String): ByteArray = Base64.decode(s, Base64.URL_SAFE or Base64.NO_WRAP)
    fun hash(s: String): String = b64(MessageDigest.getInstance("SHA-256").digest(s.toByteArray(Charsets.UTF_8)))
    fun random(): String = b64(ByteArray(32).also { SecureRandom().nextBytes(it) })
    private fun coordinate(value: BigInteger): String {
        val bytes = value.toByteArray().let { if (it.size > 32 && it[0] == 0.toByte()) it.copyOfRange(1, it.size) else it }
        require(bytes.size <= 32)
        return b64(ByteArray(32 - bytes.size) + bytes)
    }
    @Synchronized fun publicKey(alias: String): JSONObject {
        if (!store().containsAlias(alias)) {
            KeyPairGenerator.getInstance("EC", "AndroidKeyStore").apply {
                initialize(KeyGenParameterSpec.Builder(alias, KeyProperties.PURPOSE_SIGN)
                    .setAlgorithmParameterSpec(ECGenParameterSpec("secp256r1"))
                    .setDigests(KeyProperties.DIGEST_SHA256).build())
            }.generateKeyPair()
        }
        val key = store().getCertificate(alias).publicKey as ECPublicKey
        return JSONObject().put("crv", "P-256").put("kty", "EC")
            .put("x", coordinate(key.w.affineX)).put("y", coordinate(key.w.affineY))
    }
    fun thumbprint(jwk: JSONObject): String = hash("{\"crv\":\"P-256\",\"kty\":\"EC\",\"x\":\"${jwk.getString("x")}\",\"y\":\"${jwk.getString("y")}\"}")
    // Android returns ASN.1 DER; JWS ES256 requires exactly 32-byte R || 32-byte S.
    private fun joseSignature(der: ByteArray): ByteArray {
        var pos = 0
        fun byte(): Int { require(pos < der.size); return der[pos++].toInt() and 255 }
        fun length(): Int { val n = byte(); require(n < 128); return n }
        require(byte() == 0x30); require(length() == der.size - pos)
        fun integer(): ByteArray {
            require(byte() == 2)
            val n = length(); require(n in 1..33 && pos + n <= der.size)
            var bytes = der.copyOfRange(pos, pos + n); pos += n
            if (bytes.size == 33) { require(bytes[0] == 0.toByte()); bytes = bytes.copyOfRange(1, 33) }
            return ByteArray(32 - bytes.size) + bytes
        }
        val result = integer() + integer(); require(pos == der.size); return result
    }
    fun sign(alias: String, header: JSONObject, payload: JSONObject): String {
        publicKey(alias)
        val input = b64(header.toString().toByteArray(Charsets.UTF_8)) + "." + b64(payload.toString().toByteArray(Charsets.UTF_8))
        val signer = Signature.getInstance("SHA256withECDSA")
        signer.initSign(store().getKey(alias, null) as PrivateKey)
        signer.update(input.toByteArray(Charsets.US_ASCII))
        return input + "." + b64(joseSignature(signer.sign()))
    }
    fun verifyIdToken(jwt: String, jwks: JSONObject): JSONObject {
        val parts = jwt.split('.'); require(parts.size == 3) { "Malformed ID token" }
        val header = JSONObject(String(decode(parts[0]), Charsets.UTF_8))
        val alg = header.getString("alg")
        require(alg == "PS256" || alg == "RS256" || alg == "ES256") { "Unsupported ID token algorithm" }
        require(!header.has("crit")) { "Unsupported critical JWT header" }
        val kid = header.getString("kid")
        val keys = jwks.getJSONArray("keys")
        val key = (0 until keys.length()).map { keys.getJSONObject(it) }.singleOrNull {
            it.optString("kid") == kid && it.optString("kty") == (if (alg == "ES256") "EC" else "RSA") &&
                (!it.has("use") || it.optString("use") == "sig") &&
                (!it.has("alg") || it.optString("alg") == alg)
        } ?: error("ID token signing key unavailable")
        val pub: PublicKey
        val verifier: Signature
        if (alg == "ES256") {
            require(key.getString("crv") == "P-256")
            val params = AlgorithmParameters.getInstance("EC").apply { init(ECGenParameterSpec("secp256r1")) }
                .getParameterSpec(ECParameterSpec::class.java)
            pub = KeyFactory.getInstance("EC").generatePublic(ECPublicKeySpec(
                ECPoint(BigInteger(1, decode(key.getString("x"))), BigInteger(1, decode(key.getString("y")))), params))
            verifier = Signature.getInstance("SHA256withECDSA")
        } else {
            val modulus = BigInteger(1, decode(key.getString("n")))
            require(modulus.bitLength() >= 2048) { "Weak ID token signing key" }
            pub = KeyFactory.getInstance("RSA").generatePublic(RSAPublicKeySpec(modulus, BigInteger(1, decode(key.getString("e")))))
            verifier = Signature.getInstance(if (alg == "PS256") "SHA256withRSA/PSS" else "SHA256withRSA")
        }
        verifier.initVerify(pub)
        if (alg == "PS256") verifier.setParameter(PSSParameterSpec("SHA-256", "MGF1", MGF1ParameterSpec.SHA256, 32, 1))
        verifier.update((parts[0] + "." + parts[1]).toByteArray(Charsets.US_ASCII))
        require(verifier.verify(if (alg == "ES256") derSignature(decode(parts[2])) else decode(parts[2]))) { "ID token signature failed" }
        return JSONObject(String(decode(parts[1]), Charsets.UTF_8))
    }
    private fun derSignature(raw: ByteArray): ByteArray {
        require(raw.size == 64)
        fun integer(offset: Int): ByteArray {
            val v = BigInteger(1, raw.copyOfRange(offset, offset + 32)).toByteArray()
            return byteArrayOf(2, v.size.toByte()) + v
        }
        val body = integer(0) + integer(32)
        return byteArrayOf(0x30, body.size.toByte()) + body
    }
    private fun storageKey(): SecretKey {
        val alias = "setupass.pending.v1"
        if (!store().containsAlias(alias)) KeyGenerator.getInstance("AES", "AndroidKeyStore").apply {
            init(KeyGenParameterSpec.Builder(alias, KeyProperties.PURPOSE_ENCRYPT or KeyProperties.PURPOSE_DECRYPT)
                .setBlockModes(KeyProperties.BLOCK_MODE_GCM).setEncryptionPaddings(KeyProperties.ENCRYPTION_PADDING_NONE).build())
        }.generateKey()
        return store().getKey(alias, null) as SecretKey
    }
    fun encrypt(text: String): ByteArray = Cipher.getInstance("AES/GCM/NoPadding").run {
        init(Cipher.ENCRYPT_MODE, storageKey()); iv + doFinal(text.toByteArray(Charsets.UTF_8))
    }
    fun decrypt(bytes: ByteArray): String = Cipher.getInstance("AES/GCM/NoPadding").run {
        require(bytes.size > 28)
        init(Cipher.DECRYPT_MODE, storageKey(), GCMParameterSpec(128, bytes.copyOfRange(0, 12)))
        String(doFinal(bytes.copyOfRange(12, bytes.size)), Charsets.UTF_8)
    }
}
