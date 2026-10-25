. "$PSScriptRoot\Common.ps1"
Check-Keystore
Run-Docker -DockerArgs ($Compose + @('up','-d','--no-deps','esignet'))
Wait-Backend
& "$PSScriptRoot\Check.ps1"
