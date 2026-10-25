# API reference

Local issuer: `http://localhost:8088`. The UI is `http://localhost:3000`. These addresses are local demonstration endpoints, not a hosted public service. USB `adb reverse` connects the phone to the laptop's services.

## Endpoints used

| Method and path | Purpose | Authentication / format |
|---|---|---|
| GET `/.well-known/openid-configuration` | Discover issuer, token and authorization endpoints | Public JSON |
| GET `/oauth2/jwks` | Public keys for signature verification | Public JWK set |
| POST `/client-mgmt/client` | Register local parent/child client | JSON request envelope; demonstration registration policy |
| POST `/oauth2/par` | Submit parent authorization request | Form encoded; private_key_jwt client assertion and DPoP |
| GET `/oauth2/authorize` | Browser login using pushed request | `client_id` and `request_uri`; OTP/consent UI |
| POST `/oauth2/token` | Parent code exchange or child native exchange | Form encoded; private_key_jwt and fresh DPoP proof |

Client registration is performed automatically in each Android app. Its JSON envelope is `{ "requestTime": "UTC timestamp", "request": { ... } }`. The request contains `clientId`, localized client name, relying party ID, logo, redirect URIs, public JWK, claims, ACR values, grant types, authentication methods and additional configuration. The full executable examples are `SetuLogin.register()` and `BenefitsLogin.register()` in the Android source. A successful response has `response.clientId` and `response.status`. The reserved family-approval key cannot be supplied through this public request; an operator uses `deploy/pair-apps.sql`.

## Parent authorization

The PAR request includes `client_id`, `redirect_uri`, `response_type=code`, `scope=openid profile email device_sso`, state, nonce, an S256 PKCE challenge, client assertion and device proof binding. A successful response contains `request_uri` and `expires_in`. The app opens the authorization URL in the external browser. On return it verifies state and issuer and sends the authorization code with its original PKCE verifier to the token endpoint.

The parent token response contains `access_token`, `token_type=DPoP`, `expires_in`, `scope`, `id_token` and `device_secret`. The ID token includes `sid`, `ds_hash` and the encrypted device secret. The app verifies the signature and claims before persisting the session. Tokens and secrets must never be committed to Git or copied into a public demonstration.

## Native child exchange

Request body below is schematic; values must be URL encoded. Placeholder credentials are not executable tokens.

```text
grant_type=urn:ietf:params:oauth:grant-type:token-exchange
client_id=<registered-child-client>
client_assertion_type=urn:ietf:params:oauth:client-assertion-type:jwt-bearer
client_assertion=<fresh-child-signed-assertion>
subject_token_type=urn:ietf:params:oauth:token-type:id_token
subject_token=<verified-parent-id-token>
actor_token_type=urn:openid:params:token-type:device-secret
actor_token=<opaque-parent-device-secret>
scope=openid profile email
```

Header: `Content-Type: application/x-www-form-urlencoded`; `DPoP: <fresh-parent-key-proof-for-POST-token-endpoint>`. The proof is obtained from the restricted parent ContentProvider; the child does not export the parent private key. Both apps handle a server-requested DPoP nonce with one retry. The client assertion uses the client's own key and issuer/audience values implemented in the app.

A successful exchange returns a child access token and ID token. The ID token audience identifies the child and its subject/session match the parent. Token lifetime is at most 120 seconds and is bounded by the remaining parent validity. No new device secret is returned to the child.

Invalid client authentication, missing/invalid DPoP, invalid secret or ID-token binding, expired or absent session, unapproved family and unauthorized scope requests are rejected with an OAuth error response. See the provider/handler tests for the precise rejection cases; this document does not promise one identical error code for every layer.

## Android IPC API

Authority: `com.setupass.identity.sso`. Signature permission: `com.setupass.identity.NATIVE_SSO`. The provider additionally checks the caller's Binder UID, exact Benefits package and matching signing certificate.

- `credentials`: returns the valid parent ID token, opaque device secret and subject to the approved sibling.
- `proof`: produces a fresh DPoP proof only for POST to the fixed local token endpoint, with a bounded optional server nonce.
- Arbitrary URLs, arbitrary signing and generic database operations are not exposed.

## Test identity API

Mock service: POST `http://localhost:8082/v1/mock-identity-system/identity`. JSON example: `deploy/test-identity.json`; executable timestamp handling: `deploy/Seed-TestIdentity.ps1`. This is a mock service used by the prototype, not a new SetuPass identity enrollment API. Health: GET `/v1/mock-identity-system/actuator/health`.

Import `SetuPass.postman_collection.json` for public checks and the exchange reference. Perform the real success case in the Android apps; do not export Android private keys to Postman.
