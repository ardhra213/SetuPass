$ErrorActionPreference='Stop'
$repo=Split-Path $PSScriptRoot -Parent
$expected='https://github.com/ardhra213/SetuPass.git'
function Run-Git {
    param([Parameter(Mandatory=$true)][string[]]$GitArgs)
    & git.exe @GitArgs
    if ($LASTEXITCODE -ne 0) { throw 'Git stopped. Keep this window open and share the error. No force push is used.' }
}
Push-Location $repo
try {
    if (-not (Test-Path 'README.md') -or -not (Test-Path 'deploy/compose.yaml') -or -not (Test-Path 'android/settings.gradle.kts')) { throw 'Run this only inside the extracted submission package.' }
    if (-not (Test-Path '.git')) {
        $ancestor=Split-Path $repo -Parent
        while ($ancestor) {
            if (Test-Path (Join-Path $ancestor '.git')) { throw 'This folder is inside another Git repository. Extract the submission into a separate folder.' }
            $ancestor=Split-Path $ancestor -Parent
        }
        Run-Git -GitArgs @('init','-b','main')
    }
    $branch=& git.exe branch --show-current
    if ($branch -ne 'main') { throw 'Expected the submission main branch. No branch was changed.' }
    $remotes=@(& git.exe remote)
    if ($LASTEXITCODE -ne 0) { throw 'Cannot inspect repository remotes.' }
    if ($remotes -notcontains 'origin') { Run-Git -GitArgs @('remote','add','origin',$expected) }
    else {
        $remote=& git.exe remote get-url origin
        if ($LASTEXITCODE -ne 0 -or $remote -ne $expected) { throw 'Existing remote differs from the requested SetuPass repository. Nothing was pushed.' }
    }
    # Repository-local commit attribution; no changes to global Git identity.
    # This is the conventional GitHub username-based privacy address, not a verified personal email.
    Run-Git -GitArgs @('config','--local','user.name','ardhra213')
    Run-Git -GitArgs @('config','--local','user.email','ardhra213@users.noreply.github.com')
    Run-Git -GitArgs @('add','--all')
    & git.exe diff --cached --quiet
    $difference=$LASTEXITCODE
    if ($difference -eq 1) {
        Run-Git -GitArgs @('commit','-m','Import SetuPass native SSO prototype and documented acceptance evidence')
    } elseif ($difference -ne 0) { throw 'Could not inspect staged changes.' }
    Write-Host 'Uploading to https://github.com/ardhra213/SetuPass.'
    Write-Host 'If GitHub opens in your browser, sign in to ardhra213. Never paste a password or token into chat.'
    Run-Git -GitArgs @('push','-u','origin','main')
    Write-Host 'Upload complete: https://github.com/ardhra213/SetuPass' -ForegroundColor Green
} finally { Pop-Location }
