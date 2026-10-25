# Reviewer setup — separate fresh environment

This procedure is for a **new reviewer machine**, not the team's recovered laptop. The native flow passed on the team's existing environment; this fresh Compose arrangement has been statically checked but has not yet had an independent full-stack run. PostgreSQL is set to major version 17 with its matching data-directory mount. Other upstream image tags remain mutable, so future upstream changes may affect reproduction.

## Requirements

Windows with Docker Desktop's Linux engine; Android Studio and SDK 37; Git; internet for container, Go and Gradle downloads; one USB-debugging Android phone (API 24+). The Android project uses Gradle 9.6.0, Android plugin 9.4.0 and daemon JDK 25. Use a separate phone or matching debug signing certificate when the package is already installed. Do not uninstall an existing user's app to resolve certificate conflicts.

Clone into a writable directory without spaces, for example `C:\Users\reviewer\SetuPass`. Run commands from the repository root in PowerShell. All services bind to localhost; ports 3000, 8082, 8088 and 5455 must be free. Do not stop/delete another installation to free them without its owner's direction.

## Backend

```powershell
docker compose -f deploy/compose.yaml config --quiet
docker compose -f deploy/compose.yaml up -d --build
docker compose -f deploy/compose.yaml ps
curl.exe --fail http://localhost:8088/.well-known/openid-configuration
curl.exe --fail http://localhost:8082/v1/mock-identity-system/actuator/health
```

Wait for startup and an overall mock health status of UP. The Compose project is `setupass-review`. Named volumes persist PostgreSQL and both services' `/opt/mosip` directories. `init.sql` is an upstream fresh-database initializer; never run it manually against an existing database. Never use `down -v` as troubleshooting. Back up database and matching keystore together before upgrades.

Seed synthetic test data (only after mock service is healthy):

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File deploy/Seed-TestIdentity.ps1
```

Test UIN **1774231323**, test OTP **111111**. These are mock credentials, not a real identity. The seed script accepts only the specific duplicate-ID case on rerun and does not overwrite it.

## Apps

Open the repository's `android` folder in Android Studio and let Gradle sync. This creates the machine-local SDK configuration. Build/install both modules using the same local signing key. With the phone authorized and selected, either run each module in Android Studio or use:

```powershell
$adb = "$env:LOCALAPPDATA\Android\Sdk\platform-tools\adb.exe"
Push-Location android
.\gradlew.bat --no-daemon --max-workers=2 :app:assembleDebug :benefits:assembleDebug
# Continue only after BUILD SUCCESSFUL.
& $adb install -r app/build/outputs/apk/debug/app-debug.apk
& $adb install -r benefits/build/outputs/apk/debug/benefits-debug.apk
Pop-Location
& $adb reverse tcp:3000 tcp:3000
& $adb reverse tcp:8088 tcp:8088
```

Check each install returns Success. Open Setu ID and Setu Benefits once and wait for each to show Ready. This registers their distinct client keys.

## Operator approval

The pairing transaction requires exactly one active native parent and one active native child. It refuses ambiguous records.

```powershell
docker compose -f deploy/compose.yaml cp deploy/pair-apps.sql database:/tmp/pair-apps.sql
docker compose -f deploy/compose.yaml exec -T database psql -U postgres -d mosip_esignet -v ON_ERROR_STOP=1 -f /tmp/pair-apps.sql
# Continue only after COMMIT.
docker compose -f deploy/compose.yaml restart esignet
```

Wait for discovery again. Sign in to Setu ID with the test UIN/OTP, allow claims, then immediately open Benefits and tap Continue with Setu ID. It should show its signed-in state without another browser/OTP. Parent and child show the same identity reference.

## Later starts

Use `docker compose -f deploy/compose.yaml up -d`, reconnect USB and repeat both adb reverse commands. Reauthenticate Setu ID after server restart or native-session expiry. No re-pairing is needed unless client enrollment changes. Do not run the root numbered recovery/update scripts against this reviewer Compose project.
