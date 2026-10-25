# SetuPass: native SSO in eSignet Go

SetuPass demonstrates one OTP-and-consent login in **Setu ID**, followed by a distinct OAuth token issuance for **Setu Benefits** on the same Android phone. Benefits does not use the browser login session for this exchange.

## Components

| Component | Responsibility |
|---|---|
| Setu ID (`com.setupass.identity`) | PAR, PKCE, external-browser login, client assertion, DPoP, verified ID token, encrypted local session |
| Signature-protected ContentProvider | Exact sibling package / Binder UID / signing certificate checks; shares ID token and encrypted device secret; signs only fresh POST proofs for the fixed token endpoint |
| Setu Benefits (`com.setupass.benefits`) | Separate Keystore client key, private_key_jwt authentication, RFC 8693 exchange and child ID token validation |
| ThunderID fork | Public DeviceSecretProvider contract, authorization-code hook, native token-exchange branch, response and ID-claim extensions |
| eSignet provider | Secret signing and authenticated encryption, live-session map, client family and scope checks |
| PostgreSQL and persistent software keystore | Existing client/consent data and key material survive image replacement |

## First login

1. Setu ID registers a local demonstration client using a new on-device client key. `allowed_authorization_scopes` contains `device_sso`.
2. The laptop operator pairs the two exact client records. HTTP registration cannot set the reserved `native_sso_family` key. The pair script rejects zero or multiple matching clients rather than guessing.
3. Setu ID sends a PAR with `openid profile email device_sso`, state, nonce, PKCE challenge and DPoP key thumbprint.
4. eSignet UI performs mock-identity OTP and consent. The Android app validates callback state and exchanges the code using the verifier and client assertion.
5. The code handler invokes the injected DeviceSecretProvider when the authorized grant includes `device_sso`. The provider checks current enrollment, `device_sso` permission and the verified DPoP thumbprint.
6. eSignet returns ordinary tokens plus an opaque `device_secret`. The ID token includes `sid`, `ds_hash`, and the encrypted `device_secret` to satisfy the challenge's embedding requirement.
7. The Android app verifies the ID-token signature and issuer/audience/expiry/nonce, then checks the secret hash and persists the bundle encrypted with Android Keystore AES-GCM. Private EC keys never leave Keystore.

## Second-app exchange

Benefits calls the parent provider with its own Binder identity. The provider requires the exact Benefits package and the same signing certificate, in addition to Android's signature permission. There is no generic caller-supplied signing operation: proof method and URL are fixed.

Benefits authenticates itself with its own client assertion and sends:

```text
grant_type=urn:ietf:params:oauth:grant-type:token-exchange
subject_token_type=urn:ietf:params:oauth:token-type:id_token
actor_token_type=urn:openid:params:token-type:device-secret
scope=openid profile email
```

The token endpoint verifies the fresh parent-key DPoP proof and existing replay controls. The provider decrypts and verifies the device-secret JWT, verifies the signed parent ID token, and checks `sub`, parent `aud`, `sid`, `ds_hash`, expiry, live-session membership, current parent permission, same approved client family and matching `cnf.jkt`. Scope widening and requesting `device_sso` in the child fail. Audience/resource retargeting is not supported by this profile.

The child receives its own DPoP-bound access token and ID token, with the child's audience and the original subject/session. It does not receive a newly issued device secret. The access token retains the parent attribute-cache reference for downstream integration, but this app currently demonstrates authenticated access with the ID token and does not call UserInfo or submit benefit applications.

## Device-secret representation

The inner signed JWT contains `iss`, `iat`, `exp`, `sub`, `auth_id`, `sid`, `client_id`, `scope`, `acr`, `auth_time`, `claims`, `cnf.jkt`, and the attribute-cache/token-family references. It is signed with the configured eSignet signing key through RuntimeCryptoProvider, encrypted with the existing AES-GCM crypto provider, then encoded with base64url. This is an **opaque encrypted envelope containing a JWT**, not compact JWE. `ds_hash` is base64url(SHA-256(the exact returned opaque secret)). The Native SSO draft leaves both secret construction and binding details to the authorization server.

Native sessions and their parent ID tokens last up to 900 seconds. Child tokens last at most 120 seconds and cannot outlive the parent session or ID token. In-memory session state is bounded to 10,000 entries with expired-entry cleanup during issuance. Server restart invalidates native sessions; the user signs in again.

## Trusted boundaries and limitations

This is a local hackathon prototype with test identities, loopback HTTP and debug signing. The existing local registration API is not a production enrollment service. An operator with database access approves the client family; the mobile app cannot grant itself that reserved policy through registration. Local administrator access is trusted.

Before production: authenticated administrative enrollment, TLS and verified App Links, controlled release certificates, durable shared session/revocation storage, key rotation policy, mobile security review and user-facing app-sharing consent. The prototype has no refresh-token rotation or additional-scope re-consent flow. Clearing Setu ID deletes the device's stored SSO credentials and stops new proofs; it does not revoke already issued child tokens at the server. Existing child tokens can remain valid for their remaining short lifetime. Server-side RFC 7009 revocation is not yet wired to the native session map.

## Source references

- Challenge source provided by the user: PS03, Native SSO in eSignet Go.
- eSignet base: `a76b74692b4199ef977142874556e18dc55a41b1`.
- ThunderID backend base: `64f1aa911649` (the pinned dependency of that eSignet revision).
- Native SSO draft: https://openid.net/specs/openid-connect-native-sso-1_0.html
- Token exchange: https://www.rfc-editor.org/rfc/rfc8693
- DPoP: https://www.rfc-editor.org/rfc/rfc9449
