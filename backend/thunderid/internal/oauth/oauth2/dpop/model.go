// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package dpop

// VerifyParams holds the inputs required to verify a DPoP proof.
type VerifyParams struct {
	// Proof is the raw DPoP proof JWT (the value of the DPoP HTTP header).
	Proof string
	// HTM is the HTTP method of the request the proof is bound to (e.g. "POST").
	HTM string
	// HTU is the canonical request URL (without query/fragment) the proof is bound to.
	HTU string
	// AccessToken, when non-empty, requires the proof to carry a matching `ath` claim
	// equal to base64url(SHA-256(AccessToken)). Used at resource servers.
	AccessToken string
	// ExpectedJkt, when non-empty, requires the proof's computed thumbprint to equal it.
	// Used to enforce auth-code or refresh-token bindings.
	ExpectedJkt string
}

// ProofResult is the outcome of a successful DPoP proof verification.
type ProofResult struct {
	// JKT is the SHA-256 JWK thumbprint of the proof's embedded public key.
	JKT string
	// JWK is the embedded public-key JWK exactly as it appeared in the proof header.
	JWK map[string]any
	// Alg is the JWS algorithm used to sign the proof.
	Alg string
	// Confirmed is true when ExpectedJkt was supplied and matched.
	Confirmed bool
}
