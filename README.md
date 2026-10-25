# SetuPass — Native SSO for eSignet Go

**MOSIP Decode 2026 · Team: 3 brain cells · PS03**  
Repository: https://github.com/ardhra213/SetuPass

Sign in once with OTP and consent in **Setu ID**, then open **Setu Benefits** and obtain separate OAuth tokens without another browser login or OTP. Benefits is an illustrative service; it does not assess eligibility or submit real applications.

## Demonstrated

On 25 October 2026 the team built and started the native-SSO backend on Windows/Docker Desktop, installed both Android apps, approved their client family, and showed successful Benefits sign-in on a physical phone. A screenshot was reviewed; the team reports recording the demo. The video itself has not yet been reviewed in this package.

## Start here

- **New reviewer machine:** [Reviewer setup](docs/REVIEWER_SETUP.md). This separate clean-environment procedure is supplied for reproduction and has not yet been run end-to-end. Do not run it on the team's existing installation.
- **Team's already working installation:** keep it running; no reinstall is needed for submission. The root numbered scripts are the recovered-installation update workflow, not a fresh-install workflow.
- [Architecture and limitations](docs/ARCHITECTURE.md)
- [API reference](docs/API_REFERENCE.md) and [Postman collection](docs/SetuPass.postman_collection.json)
- [Test results and evidence](docs/TEST_REPORT.md)
- [Developer guide](docs/DEVELOPER_GUIDE.md)
- [Submission checklist](docs/SUBMISSION_CHECKLIST.md)

## Contents

`android/`: Kotlin/Compose parent and child apps.  
`backend/esignet-service/`: eSignet Go host and device-secret provider.  
`backend/thunderid/`: engine fork and native token-exchange handling.  
`deploy/`: standalone local reviewer Compose, schema, proxy, test identity and family approval SQL.  
`docs/`: architecture, tests, API reference and pitch deck.

## Security boundaries

The second app must have the approved package/signing certificate and an approved server-side client family. Device proofs are created with an Android Keystore key. The backend checks the encrypted device secret, signed parent ID token, session/subject/hash binding, DPoP thumbprint, client family and scope limits. Parent native sessions last up to 15 minutes; child tokens at most 2 minutes. Restarting eSignet invalidates native sessions. Local session clearing is not server-wide revocation.

This is a local test prototype using loopback HTTP, debug signing, mock identities and demonstration database passwords. It is not production-ready or protocol-certified. See architecture for exact limitations.

## Source provenance and development history

The included eSignet host derives from upstream revision `a76b74692b4199ef977142874556e18dc55a41b1`; ThunderID derives from `64f1aa911649`. Their original licenses and notices are retained in the corresponding source directories. New implementation areas are listed in the developer guide. Development was AI-assisted and tested as documented.

This repository is a submission import of work developed locally. Its first upload does not reconstruct earlier local edit history or represent upstream code as newly authored work. Preserve any existing original Git history separately; do not backdate commits or invent team contributions. Subsequent genuine changes should be committed normally.
