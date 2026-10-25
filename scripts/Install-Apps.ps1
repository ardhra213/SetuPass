. "$PSScriptRoot\Common.ps1"
$project=Join-Path $env:USERPROFILE 'SetuPassAndroid'
$adb="$env:LOCALAPPDATA\Android\Sdk\platform-tools\adb.exe"
if (-not (Test-Path $adb)) { throw 'Android SDK not found at its usual location.' }
$devices=& $adb devices
if (-not ($devices | Select-String '\sdevice$')) { throw 'Connect and unlock your phone, enable USB debugging and accept its authorization prompt.' }
if (($devices | Select-String '\sdevice$').Count -ne 1) { throw 'Connect exactly one phone and close any running emulator.' }
if (Test-Path $project) {
    $stamp=Get-Date -Format 'yyyyMMdd-HHmmss'
    Move-Item $project "$project-before-$stamp"
}
New-Item -ItemType Directory -Path $project -Force | Out-Null
Copy-Item "$Bundle\android\*" $project -Recurse -Force
$sdk=Join-Path $env:LOCALAPPDATA 'Android\Sdk'
('sdk.dir=' + ($sdk -replace '\\','/')) | Set-Content "$project\local.properties" -Encoding ASCII
$jbr='C:\Program Files\Android\Android Studio\jbr'
if (Test-Path "$jbr\bin\java.exe") { $env:JAVA_HOME=$jbr }
Push-Location $project
try {
    Write-Host 'Building both apps with your existing Android debug signing key. The first build can take several minutes.'
    & .\gradlew.bat --no-daemon --max-workers=2 :app:assembleDebug :benefits:assembleDebug
    if ($LASTEXITCODE -ne 0) { throw 'Android build stopped. Keep this output; your original Android Studio project is unchanged.' }
    & $adb install -r '.\app\build\outputs\apk\debug\app-debug.apk'
    if ($LASTEXITCODE -ne 0) { throw 'Parent update failed. Do not uninstall it; share the install error.' }
    & $adb install -r '.\benefits\build\outputs\apk\debug\benefits-debug.apk'
    if ($LASTEXITCODE -ne 0) { throw 'Benefits installation failed.' }
    & $adb reverse tcp:3000 tcp:3000
    & $adb reverse tcp:8088 tcp:8088
    & $adb shell am start -n com.setupass.identity/.MainActivity
    Write-Host 'Let Setu ID show Ready, then open Setu Benefits from your phone and let it show Ready.' -ForegroundColor Green
    Write-Host 'Then run 03_PAIR_APPS.cmd once. You do not need to copy any IDs or tokens.'
} finally { Pop-Location }
