package com.setupass.identity

import android.content.Context
import org.json.JSONObject
import java.io.File

internal object NativeSession {
    private fun file(c: Context) = File(c.noBackupFilesDir, "native-sso.enc")
    @Synchronized fun save(c: Context, id: String, secret: String, expiry: Long, subject: String) {
        val data = JSONObject().put("id_token", id).put("device_secret", secret).put("expires", expiry).put("sub", subject)
        val target = file(c); val temp = File(c.noBackupFilesDir, "native-sso.tmp")
        temp.writeBytes(AuthCrypto.encrypt(data.toString()))
        check(temp.renameTo(target)) { "Could not store native session" }
    }
    @Synchronized fun get(c: Context): JSONObject {
        check(file(c).exists()) { "Open Setu ID and sign in first." }
        val data = JSONObject(AuthCrypto.decrypt(file(c).readBytes()))
        check(data.getLong("expires") > System.currentTimeMillis() / 1000) { "Setu ID session expired. Sign in again." }
        return data
    }
    @Synchronized fun clear(c: Context) { file(c).delete() }
}
