. "$PSScriptRoot\Common.ps1"
Check-Keystore
$stamp=Get-Date -Format 'yyyyMMdd-HHmmss'
$backup=Join-Path $BackupRoot "native-sso-$stamp"
New-Item -ItemType Directory -Path $backup -Force | Out-Null
$oldImage=(& docker inspect --format '{{.Image}}' setupass-esignet-1).Trim()
if ($LASTEXITCODE -ne 0 -or -not $oldImage) { throw 'Cannot identify the working image.' }
Run-Docker -DockerArgs @('tag',$oldImage,"setupass/esignet:before-native-$stamp")
$oldImage | Set-Content "$backup\image-id.txt" -Encoding ASCII
Copy-Item "$Root\compose.setupass.yaml" "$backup\compose.setupass.yaml"
if (Test-Path "$Root\compose.native-sso.yaml") { Copy-Item "$Root\compose.native-sso.yaml" "$backup\previous-native.yaml" }
Write-Host 'Building the new backend. Your existing backend is still running.'
Run-Docker -DockerArgs @('build','--build-arg',"BASE_IMAGE=setupass/esignet:before-native-$stamp",'-t','setupass/esignet:native-sso-v1',"$Bundle\backend")
Write-Host 'Creating database and keystore backups before switching images.'
Backup-Database "$backup\database.sql"
Run-Docker -DockerArgs @('stop','setupass-esignet-1')
try {
Copy-Item "$Root\state\esignet-keystore" "$backup\keystore" -Recurse
$hashBefore=(Get-FileHash "$Root\state\esignet-keystore\test.pfx" -Algorithm SHA256).Hash
@'
services:
  esignet:
    image: setupass/esignet:native-sso-v1
    environment:
      MOSIP_ESIGNET_SIGNING_KEY_REF_ID: RSA_2048
'@ | Set-Content "$Root\compose.native-sso.yaml" -Encoding ASCII
    Run-Docker -DockerArgs ($Compose + @('config','--quiet'))
    Run-Docker -DockerArgs ($Compose + @('up','-d','--no-deps','esignet'))
    Check-Keystore
    Wait-Backend
    Write-Host 'Backend ready. Existing clients, consents and keystore have been preserved.' -ForegroundColor Green
    Write-Host "Backup: $backup"
    Write-Host 'Next: run 02_INSTALL_APPS.cmd with your phone connected and unlocked.'
} catch {
    Write-Host 'Update failed. Returning to the previous image using the same persistent keystore.' -ForegroundColor Yellow
    "services:`n  esignet:`n    image: setupass/esignet:before-native-$stamp" | Set-Content "$Root\compose.native-sso.yaml" -Encoding ASCII
    Run-Docker -DockerArgs ($Compose + @('up','-d','--no-deps','esignet'))
    throw
}
