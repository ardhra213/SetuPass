// Copyright 2025 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package jwks

// JWKS defines the structure of a JSON Web Key Set.
type JWKS struct {
	Kid     string   `json:"kid,omitempty"`
	Kty     string   `json:"kty"`
	Use     string   `json:"use,omitempty"`
	Alg     string   `json:"alg,omitempty"`
	N       string   `json:"n,omitempty"`
	E       string   `json:"e,omitempty"`
	Crv     string   `json:"crv,omitempty"`
	X       string   `json:"x,omitempty"`
	Y       string   `json:"y,omitempty"`
	Pub     string   `json:"pub,omitempty"`
	X5c     []string `json:"x5c,omitempty"`
	X5t     string   `json:"x5t,omitempty"`
	X5tS256 string   `json:"x5t#S256,omitempty"`
}

// JWKSResponse defines the structure of the response containing JWKS.
type JWKSResponse struct {
	Keys []JWKS `json:"keys"`
}
