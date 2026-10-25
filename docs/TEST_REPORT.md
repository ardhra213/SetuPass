# Validation report

Updated 25 October 2026. Development tests are distinguished from team-provided Windows and physical-phone acceptance evidence below.

## Completed

| Check | Result | Scope |
|---|---|---|
| eSignet executable build | PASS | Go 1.26, included ThunderID fork linked locally |
| eSignet engine and client-management tests | PASS with race detector | Host provider and registration behavior |
| ThunderID OAuth2 package tree | PASS with race detector | 24 tested packages across the combined command, including the two host packages; packages without tests are separately reported |
| Native secret cryptographic tests | PASS | Real RSA-PSS signatures and AES-GCM in an isolated test provider |
| Native exchange handler test | PASS | Child token audience/identity, DPoP binding, remaining-lifetime limit, session/cache/family propagation, no new child device secret |
| Parent Android app | BUILD SUCCESSFUL | Debug APK produced with Gradle 9.6.0 / Android plugin 9.4.0 / Android SDK 37 |
| Child Android app | BUILD SUCCESSFUL | Separate package; same build invocation |
| PowerShell scripts | PASS | Parsed with PowerShell 7.5.4; Docker argument forwarding tested with a stub, including spaces |
| PowerPoint | PASS | Eight slides rendered and visually inspected; package/layout validation passed |

The full Android build completed successfully: 72 tasks executed. Build artifacts are intentionally not the installation route because this environment's debug certificate differs from the user's existing certificate. The installer rebuilds on Windows with the user's existing key.

## Native secret rejection cases exercised

Wrong or missing device proof thumbprint; widened scope; `device_sso` requested by the child; empty scopes; using the parent as the child; tampered encrypted secret; tampered ID-token signature; mismatched subject, audience, session ID, secret hash or issuer; expired ID token; a different client family; removal of the parent's `device_sso` permission; and missing live-session state.

The handler test also rejects a missing verified DPoP context, provider denial, and a resource-retargeting request. Existing upstream DPoP tests run as part of the OAuth2 suite. This does not constitute an independent security audit or a protocol certification.

## Windows and phone acceptance — 25 October 2026

Evidence supplied by the team in the build session:

- Docker build finished successfully; image `setupass/esignet:native-sso-v1` started with the recovered database and persistent keystore.
- Installer reported preserving clients, consents and keystore and saving a backup before replacement.
- Both apps installed after moving the installer/project into writable user-profile directories.
- Pairing transaction completed with BEGIN, DO and COMMIT; approved parent/child family was reported.
- Phone screenshot shows Setu Benefits signed in through Setu ID and the matching identity-reference prefix. This confirms the positive native flow UI result; no raw tokens were collected in this evidence.
- Team reports completing a demo video; video not reviewed here.
- Final check: four containers running; issuer `http://localhost:8088`; token exchange advertised `True`; 4 client rows; 3 consent rows; keystore 13,476 bytes.

The initial Windows PowerShell 5.1 argument-forwarding error was fixed with explicit named string-array parameters. The subsequent real backend build completed. The earlier PowerShell 7-only check did not catch that compatibility issue.

## Remaining verification and limits

- Fresh reviewer setup is documented but has not been executed end-to-end in an independent environment.
- Unauthorized-caller Android IPC rejection has not been demonstrated on the physical phone. Positive sibling access passed; source checks and backend negative tests do not replace that device test.
- A deliberate later container-recreation persistence test has not been run after the final successful phone demo.
- Submission upload and demo video review remain to be completed.

## Re-run

From `backend/` with Go 1.26:

```sh
go test -race ./esignet-service/internal/engine ./esignet-service/internal/clientmgmt ./thunderid/internal/oauth/oauth2/...
go build ./esignet-service/cmd/esignet
```

From `android/` with the installed Android SDK:

```powershell
.\gradlew.bat --no-daemon --max-workers=2 :app:assembleDebug :benefits:assembleDebug
```

The archive contains concise Go test output under `docs/verification/`. No live credentials are included.
