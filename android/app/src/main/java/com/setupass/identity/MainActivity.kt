package com.setupass.identity

import android.content.Intent
import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.activity.enableEdgeToEdge
import androidx.compose.foundation.background
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.rememberScrollState
import androidx.compose.foundation.verticalScroll
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import com.setupass.identity.ui.theme.SetuPassTheme

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        SetuLogin.initialize(applicationContext)
        enableEdgeToEdge()
        setContent { SetuPassTheme { SetuPassHome {
            SetuLogin.start { uri -> startActivity(Intent(Intent.ACTION_VIEW, uri).addCategory(Intent.CATEGORY_BROWSABLE)) }
        } } }
        acceptCallback(intent)
    }
    override fun onNewIntent(intent: Intent) {
        super.onNewIntent(intent)
        setIntent(intent)
        acceptCallback(intent)
    }
    private fun acceptCallback(incoming: Intent) {
        incoming.data?.let { uri ->
            if (uri.scheme == "com.setupass.identity" && uri.host == "oauth" && uri.path == "/callback") {
                incoming.data = null
                SetuLogin.receive(uri)
            }
        }
    }
}

@Composable
fun SetuPassHome(onLogin: () -> Unit) {
    val navy = Color(0xFF10243A)
    val teal = Color(0xFF007F78)
    Column(
        modifier = Modifier.fillMaxSize().background(Color(0xFFF4F7FA)).safeDrawingPadding()
            .verticalScroll(rememberScrollState()).padding(24.dp),
        verticalArrangement = Arrangement.spacedBy(20.dp)
    ) {
        Spacer(Modifier.height(16.dp))
        Text("SETUPASS", color = teal, fontWeight = FontWeight.Bold, style = MaterialTheme.typography.labelLarge)
        Text("Your identity.\nYour permission.", color = navy, fontWeight = FontWeight.Bold, style = MaterialTheme.typography.headlineLarge)
        Text("Sign in once to begin your Setu ID journey.", color = Color(0xFF526477), style = MaterialTheme.typography.bodyLarge)
        Card(colors = CardDefaults.cardColors(containerColor = Color.White), modifier = Modifier.fillMaxWidth()) {
            Column(Modifier.padding(24.dp), verticalArrangement = Arrangement.spacedBy(16.dp)) {
                Text("Setu ID", color = navy, fontWeight = FontWeight.Bold, style = MaterialTheme.typography.titleLarge)
                Text(if (SetuLogin.signedIn) "You're signed in" else "Welcome", style = MaterialTheme.typography.titleMedium)
                Text(SetuLogin.status, style = MaterialTheme.typography.bodyMedium)
                if (SetuLogin.busy) LinearProgressIndicator(modifier = Modifier.fillMaxWidth())
                if (SetuLogin.signedIn) {
                    Text("Identity reference: ${SetuLogin.subject.take(12)}…", style = MaterialTheme.typography.labelMedium)
                    OutlinedButton(onClick = { SetuLogin.clearSession() }, modifier = Modifier.fillMaxWidth()) { Text("Clear app session") }
                } else {
                    Button(onClick = onLogin, enabled = !SetuLogin.busy,
                        colors = ButtonDefaults.buttonColors(containerColor = teal), modifier = Modifier.fillMaxWidth()) {
                        Text(if (SetuLogin.busy) "Please wait…" else "Sign in with eSignet")
                    }
                }
            }
        }
        Text("Your browser will ask for your test identity, OTP and consent, then return you here.", color = Color(0xFF526477))
        Text("Local hackathon prototype • Test identities only", style = MaterialTheme.typography.labelSmall, color = Color(0xFF526477))
    }
}
