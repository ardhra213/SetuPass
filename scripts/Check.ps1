. "$PSScriptRoot\Common.ps1"
Check-Keystore
Run-Docker -DockerArgs @('ps','-a','--format','table {{.Names}}\t{{.Status}}')
$r=Invoke-RestMethod 'http://localhost:8088/.well-known/openid-configuration' -TimeoutSec 10
Write-Host "Issuer: $($r.issuer)"
Write-Host "Token exchange advertised: $($r.grant_types_supported -contains 'urn:ietf:params:oauth:grant-type:token-exchange')"
Run-Docker -DockerArgs @('exec','setupass-database-1','psql','-U','postgres','-d','mosip_esignet','-c',"SELECT 'clients' AS kind,COUNT(*) FROM esignet.client_detail UNION ALL SELECT 'consents',COUNT(*) FROM esignet.consent_detail;")
Write-Host "Keystore bytes: $((Get-Item "$Root\state\esignet-keystore\test.pfx").Length)"
$adb="$env:LOCALAPPDATA\Android\Sdk\platform-tools\adb.exe"
if(Test-Path $adb){ & $adb reverse tcp:3000 tcp:3000; & $adb reverse tcp:8088 tcp:8088 }
Write-Host 'No tokens or private key contents were printed.'
