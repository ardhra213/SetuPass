# Developer guide

## Reproducible source

The archive includes complete eSignet service and ThunderID backend source at the base revisions identified in ARCHITECTURE.md, plus the new native-SSO implementation. The original repository licenses and notices are retained. No user database backup or private runtime keystore is included. Upstream test fixtures remain in the source.

Go 1.26 is required. From `backend/`, the `go.work` file links the host to the included fork without changing the pinned upstream go.mod:

```sh
go test -race ./esignet-service/internal/engine ./esignet-service/internal/clientmgmt ./thunderid/internal/oauth/oauth2/...
go build ./esignet-service/cmd/esignet
```

The Docker build uses an isolated local module replacement inside its build stage. `BASE_IMAGE` defaults to the upstream image for developers; the Windows update script overrides it with a tag of the currently working local image. This preserves the runtime environment during the update. It uses the existing Compose files and adds only an image/signing-key override.

## Main changes

- `thunderid/pkg/thunderidengine/providers/device_secret.go`: public provider contract and exchanged state.
- `thunderid/internal/oauth/oauth2/authz`: retains authentication request ID on authorization codes.
- `thunderid/internal/oauth/oauth2/granthandlers/authorization_code.go`: device-secret issuance and parent token claims.
- `thunderid/internal/oauth/oauth2/granthandlers/native_sso.go`: native exchange profile.
- `thunderid/internal/oauth/oauth2/token/service.go`: rejects disallowed child scopes before scope filtering; returns device secret only where issued.
- `thunderid/internal/oauth/oauth2/tokenservice`: bounded ID token lifetime and native claims.
- `esignet-service/internal/engine/device_secret_provider.go`: cryptographic and policy validation.
- `esignet-service/internal/engine/actor_provider.go`: supplies the provider and registered grants.
- `esignet-service/internal/clientmgmt/validate.go`: recognizes token-exchange registration.
- `android/app`: parent login and restricted IPC.
- `android/benefits`: separate client and token exchange.

## Android build

The supplied project keeps the user's working Gradle/Android plugin versions. It requires the Android SDK and JDK requested by the Gradle daemon configuration. Both modules use the same local debug signing key by default. Build with:

```powershell
.\gradlew.bat --no-daemon --max-workers=2 :app:assembleDebug :benefits:assembleDebug
```

The Windows installer builds on the user's machine so the updated parent retains the original signing identity. APKs generated on another machine normally use a different debug certificate and cannot update that installed app. Do not solve a certificate mismatch by uninstalling without understanding the loss of Android Keystore credentials.

## Local family approval

The public registration validator rejects `native_sso_family`; this key is absent from its allowed-key list. The local administrator script adds that policy directly in PostgreSQL only after one parent and one child record exist, and restarts eSignet to invalidate cached client configuration. This is a demonstration enrollment workflow, not an API for production.

Repeated pairing after login invalidates in-memory sessions. Sign in again after pairing or restarting the backend. Do not rename client records or erase the existing two historical clients/consents to make the pair script pass.

## Postman

Import SetuPass.postman_collection.json. The first three requests check discovery, public keys and rejection of unauthenticated token exchange. The fourth is a reference for the protocol and is skipped without a fresh proof. The real successful exchange runs inside Android so the device key never needs to be exported. The Go tests provide repeatable negative binding/scope coverage without copying secrets into Postman.

## Recovery

The backend update creates a timestamped directory under the existing SetuPass_Backups folder. It contains the old image ID, database dump, Compose snapshot and keystore copy. The failed-startup path switches back to the captured image tag with the same persistent keystore. No automatic database restore or deletion occurs.

For a later manual rollback, find the tag recorded by the update, replace the image value in `compose.native-sso.yaml`, and run `05_START_AGAIN.cmd`. Retain `compose.setupass.yaml` and its `/opt/mosip` bind mount. Never restore only a database alias table without its matching keys.
