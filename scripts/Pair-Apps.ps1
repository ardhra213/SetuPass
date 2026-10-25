. "$PSScriptRoot\Common.ps1"
Check-Keystore
$stamp=Get-Date -Format 'yyyyMMdd-HHmmss'
New-Item -ItemType Directory -Force $BackupRoot | Out-Null
Backup-Database "$BackupRoot\before-app-pair-$stamp.sql"
# Exact cardinality avoids silently approving an old or ambiguous client. Reserved config key is not accepted by the HTTP registration API.
$sql=@'
BEGIN;
DO $pair$
DECLARE parent_id text; child_id text; n integer;
BEGIN
 SELECT COUNT(*) INTO n FROM esignet.client_detail WHERE id LIKE 'setupass-native-parent-%' AND status='ACTIVE';
 IF n<>1 THEN RAISE EXCEPTION 'Expected one native parent. Open Setu ID and wait for Ready. Found %',n; END IF;
 SELECT COUNT(*) INTO n FROM esignet.client_detail WHERE id LIKE 'setupass-native-child-%' AND status='ACTIVE';
 IF n<>1 THEN RAISE EXCEPTION 'Expected one native child. Open Setu Benefits and wait for Ready. Found %',n; END IF;
 SELECT id INTO parent_id FROM esignet.client_detail WHERE id LIKE 'setupass-native-parent-%' AND status='ACTIVE';
 SELECT id INTO child_id FROM esignet.client_detail WHERE id LIKE 'setupass-native-child-%' AND status='ACTIVE';
 UPDATE esignet.client_detail SET additional_config=(COALESCE(NULLIF(additional_config,''),'{}')::jsonb || jsonb_build_object('native_sso_family','setupass-local-approved'))::text WHERE id IN (parent_id,child_id);
 RAISE NOTICE 'Approved Setu ID and Setu Benefits for local native SSO';
END $pair$;
COMMIT;
'@
$tmp=Join-Path $env:TEMP 'setupass-pair.sql'
[IO.File]::WriteAllText($tmp,$sql,[Text.UTF8Encoding]::new($false))
Run-Docker -DockerArgs @('cp',$tmp,'setupass-database-1:/tmp/setupass-pair.sql')
Run-Docker -DockerArgs @('exec','setupass-database-1','psql','-U','postgres','-d','mosip_esignet','-v','ON_ERROR_STOP=1','-f','/tmp/setupass-pair.sql')
Run-Docker -DockerArgs @('restart','setupass-esignet-1')
Wait-Backend
Write-Host 'Paired. Sign in to Setu ID with UIN 1774231323 and test OTP 111111. Allow the claims, then open Setu Benefits and tap Continue with Setu ID.' -ForegroundColor Green
