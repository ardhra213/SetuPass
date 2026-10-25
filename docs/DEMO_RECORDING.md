# Real phone demo: about 3 minutes

Only record after both apps work. Use the phone's screen recorder. Hide notifications, keep USB connected, and use only the mock identity. Do not display JWTs, private keys, database dumps or keystore contents.

| Time | Screen/action | Suggested narration |
|---|---|---|
| 0:00–0:20 | Setu ID home | “SetuPass demonstrates native single sign-on between two Android apps using eSignet Go. I will sign in once and then enter a second app without another OTP.” |
| 0:20–1:05 | Tap Sign in, complete test OTP and consent | “The first login uses the system browser, PKCE, pushed authorization requests and a device-bound proof. The user approves the requested claims.” |
| 1:05–1:25 | Parent signed-in card | “Setu ID validates the server's ID token and stores an encrypted device-session bundle. Its private key stays in Android Keystore.” |
| 1:25–2:00 | Open Setu Benefits and tap Continue with Setu ID | “This is a separate app with its own OAuth client. It obtains session material through a certificate-protected Android interface and exchanges it for its own tokens. There is no second browser login or OTP.” |
| 2:00–2:20 | Compare identity-reference prefixes | “Both apps refer to the same authenticated person. The second app receives tokens intended for its own client.” |
| 2:20–2:40 | Show security-test report on laptop, if recording allows | “Automated tests reject tampering, the wrong proof key, another client family, expired or mismatched claims, and wider scopes.” |
| 2:40–3:00 | Pitch closing | “The demo uses a mock identity system and local HTTP. Production work includes authenticated enrollment, durable revocation and release-signing controls.” |

Optional negative demonstration: clear Setu ID's app session, then clear Benefits' session and tap Continue again. Benefits should ask for parent sign-in. This demonstrates the local broker behavior, not immediate server revocation of an already issued access token.

Do not claim actual benefit eligibility, Aadhaar integration, production readiness, or a measured login-time reduction. The Benefits service cards illustrate a possible application of SSO.
