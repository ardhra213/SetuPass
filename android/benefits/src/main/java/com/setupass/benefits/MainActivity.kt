package com.setupass.benefits

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
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp

class MainActivity: ComponentActivity() {
    override fun onCreate(savedInstanceState:Bundle?) {
        super.onCreate(savedInstanceState);enableEdgeToEdge();BenefitsLogin.initialize(applicationContext)
        setContent { MaterialTheme {
            Column(Modifier.fillMaxSize().background(Color(0xFFF4F7FA)).safeDrawingPadding().verticalScroll(rememberScrollState()).padding(24.dp),verticalArrangement=Arrangement.spacedBy(20.dp)) {
                Text("SETUPASS / BENEFITS",color=Color(0xFF007F78),fontWeight=FontWeight.Bold)
                Text("One sign-in.\nMore possibilities.",style=MaterialTheme.typography.headlineLarge,fontWeight=FontWeight.Bold,color=Color(0xFF10243A))
                Text("Continue using your verified Setu ID session on this phone.")
                Card(Modifier.fillMaxWidth(),colors=CardDefaults.cardColors(containerColor=Color.White)) {
                    Column(Modifier.padding(24.dp),verticalArrangement=Arrangement.spacedBy(16.dp)) {
                        Text(if(BenefitsLogin.signedIn)"Welcome to Setu Benefits" else "Your identity, carried forward",style=MaterialTheme.typography.titleLarge)
                        Text(BenefitsLogin.status)
                        if(BenefitsLogin.busy)LinearProgressIndicator(Modifier.fillMaxWidth())
                        if(BenefitsLogin.signedIn) {
                            Text("Identity reference: ${BenefitsLogin.subject.take(12)}…")
                            Text("Demo service access unlocked",fontWeight=FontWeight.Bold)
                            Text("Explore education support, document help and community services. These are illustrative services; no eligibility decision or application is submitted.")
                            OutlinedButton(onClick={BenefitsLogin.clear()}){Text("Clear Benefits session")}
                        } else Button(onClick={BenefitsLogin.continueSSO()},enabled=!BenefitsLogin.busy,modifier=Modifier.fillMaxWidth()){Text("Continue with Setu ID")}
                    }
                }
                OutlinedButton(onClick={
                    packageManager.getLaunchIntentForPackage("com.setupass.identity")?.let{startActivity(it)}
                }){Text("Open Setu ID")}
                Text("Local hackathon prototype • Test identities only",style=MaterialTheme.typography.labelSmall)
            }
        } }
    }
}
