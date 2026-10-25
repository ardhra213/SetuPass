package com.setupass.identity

import android.content.ContentProvider
import android.content.ContentValues
import android.content.pm.PackageManager
import android.database.Cursor
import android.net.Uri
import android.os.Binder
import android.os.Bundle

/** Only the explicitly named sibling with the same signing certificate can use this broker. */
class SsoProvider : ContentProvider() {
    override fun onCreate() = true
    override fun call(method: String, arg: String?, extras: Bundle?): Bundle {
        val c = requireNotNull(context)
        val uid = Binder.getCallingUid()
        val packages = c.packageManager.getPackagesForUid(uid)?.toList() ?: emptyList()
        check(packages == listOf("com.setupass.benefits") &&
            c.packageManager.checkSignatures(c.applicationInfo.uid, uid) == PackageManager.SIGNATURE_MATCH) { "Untrusted app" }
        return try {
            val state = NativeSession.get(c)
            when (method) {
                "credentials" -> Bundle().apply {
                    putString("id_token", state.getString("id_token"))
                    putString("device_secret", state.getString("device_secret"))
                    putString("subject", state.getString("sub"))
                }
                "proof" -> Bundle().apply {
                    // No caller-selected URL, verb or signing payload: this is not a generic signing oracle.
                    val nonce = extras?.getString("nonce")
                    require(nonce == null || nonce.length <= 512)
                    putString("dpop", SetuLogin.proof("http://localhost:8088/oauth2/token", "POST", nonce = nonce))
                }
                else -> error("Unsupported broker operation")
            }
        } catch (_: Exception) { Bundle().apply { putString("error", "Open Setu ID and sign in again.") } }
    }
    override fun query(u: Uri, p: Array<out String>?, s: String?, a: Array<out String>?, o: String?): Cursor? = null
    override fun getType(u: Uri): String? = null
    override fun insert(u: Uri, v: ContentValues?): Uri? = throw UnsupportedOperationException()
    override fun delete(u: Uri, s: String?, a: Array<out String>?): Int = throw UnsupportedOperationException()
    override fun update(u: Uri, v: ContentValues?, s: String?, a: Array<out String>?): Int = throw UnsupportedOperationException()
}
