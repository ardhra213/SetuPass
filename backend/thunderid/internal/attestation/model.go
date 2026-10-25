// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package attestation

// attestationObject is the CBOR-encoded structure produced by DCAppAttestService.attestKey.
type attestationObject struct {
	Fmt     string `cbor:"fmt"`
	AttStmt struct {
		X5C [][]byte `cbor:"x5c"`
		// Receipt is an App Store receipt unrelated to binary identity; intentionally unused.
		Receipt []byte `cbor:"receipt"`
	} `cbor:"attStmt"`
	AuthData []byte `cbor:"authData"`
}

// parsedAuthData holds the authenticator data fields relevant to attestation verification.
type parsedAuthData struct {
	rpIDHash     []byte
	signCount    uint32
	aaguid       [16]byte
	credentialID []byte
}
