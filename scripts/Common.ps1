$ErrorActionPreference = 'Stop'
$Root = $env:SETUPASS_EXISTING_ROOT
if (-not $Root) { $Root = Join-Path $env:USERPROFILE 'SetuPass_Windows_Starter' }
$Bundle = Split-Path $PSScriptRoot -Parent
$BackupRoot = Join-Path (Split-Path $Root -Parent) 'SetuPass_Backups'
$Compose = @('compose','-p','setupass','-f',"$Root\source\esignet\docker-compose\docker-compose.yaml",'-f',"$Root\compose.setupass.yaml",'-f',"$Root\compose.native-sso.yaml")
function Run-Docker {
    param([Parameter(Mandatory=$true)][string[]]$DockerArgs)
    & docker.exe @DockerArgs
    if ($LASTEXITCODE -ne 0) { throw "Docker command failed. Keep this window open and share the last error." }
}
function Check-Keystore {
    $inspect = (& docker inspect setupass-esignet-1 | ConvertFrom-Json)[0]
    if ($LASTEXITCODE -ne 0) { throw 'Cannot inspect the existing backend.' }
    if (-not ($inspect.Mounts | Where-Object { $_.Destination -eq '/opt/mosip' -and $_.Type -eq 'bind' })) {
        throw 'STOP: The recovered /opt/mosip bind mount is missing. No container was replaced.'
    }
    $key = "$Root\state\esignet-keystore\test.pfx"
    if (-not (Test-Path $key) -or (Get-Item $key).Length -lt 1000) {
        throw 'STOP: The recovered keystore is missing or unexpectedly small. No keys will be reset.'
    }
}
function Backup-Database([string]$Destination) {
    Run-Docker -DockerArgs @('exec','setupass-database-1','pg_dumpall','-U','postgres','-f','/tmp/setupass-backup.sql')
    Run-Docker -DockerArgs @('cp','setupass-database-1:/tmp/setupass-backup.sql', $Destination)
    if ((Get-Item $Destination).Length -lt 1000) { throw 'Database backup is unexpectedly small.' }
}
function Wait-Backend {
    for ($i=0; $i -lt 45; $i++) {
        try {
            $r = Invoke-RestMethod -Uri 'http://localhost:8088/.well-known/openid-configuration' -TimeoutSec 3
            if ($r.issuer -eq 'http://localhost:8088') { return }
        } catch { }
        Start-Sleep -Seconds 2
    }
    throw 'The backend did not become available. Run 04_CHECK.cmd and share its output.'
}
