$ErrorActionPreference='Stop'
$body=Get-Content "$PSScriptRoot\test-identity.json" -Raw | ConvertFrom-Json
$body.requestTime=[DateTime]::UtcNow.ToString('yyyy-MM-ddTHH:mm:ss.fffZ')
$r=Invoke-RestMethod -Method Post -Uri 'http://localhost:8082/v1/mock-identity-system/identity' -ContentType 'application/json' -Body ($body | ConvertTo-Json -Depth 20)
if ($r.errors) {
    if (@($r.errors | Where-Object { $_.errorCode -ne 'duplicate_individual_id' }).Count -gt 0) { throw ($r.errors | ConvertTo-Json -Compress) }
    Write-Host 'Test identity already exists; it was not changed.'
} else { Write-Host 'Synthetic test identity created. UIN 1774231323; test OTP 111111.' }
