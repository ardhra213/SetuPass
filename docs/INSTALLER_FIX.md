# Windows PowerShell argument fix

The original installer used ValueFromRemainingArguments, which treats an array as one element in Windows PowerShell 5.1. Docker then received the entire tag command as one subcommand.

The corrected scripts use an explicit named string-array parameter, -DockerArgs, and invoke docker.exe with array splatting. All wrapper call sites were updated, including backup, pairing, checking and startup. The original update stopped at tagging, before backend stop or data changes.

Validation for this correction: archive integrity checked; every Docker wrapper call inspected for the explicit argument parameter; no ValueFromRemainingArguments remains. Windows PowerShell 5.1 execution is not available in this development environment. The original PowerShell 7 checks did not catch this version difference.

Reference: https://learn.microsoft.com/en-us/powershell/module/microsoft.powershell.core/about/about_functions_advanced_parameters?view=powershell-5.1
