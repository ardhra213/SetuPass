// SPDX-License-Identifier: MPL-2.0
package engine

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/mosip/esignet/internal/config"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

// Exercise real signature and authenticated-encryption operations while isolating database/key plumbing.
type nativeTestCrypto struct {
	providers.RuntimeCryptoProvider
	key *rsa.PrivateKey
	aes []byte
}

func (c nativeTestCrypto) Sign(_ context.Context, _ providers.KeyRef, _ string, b []byte) ([]byte, error) {
	h := sha256.Sum256(b)
	return rsa.SignPSS(rand.Reader, c.key, crypto.SHA256, h[:], &rsa.PSSOptions{SaltLength: 32})
}
func (c nativeTestCrypto) Verify(_ context.Context, _ providers.KeyRef, _ string, b, s []byte) error {
	h := sha256.Sum256(b)
	return rsa.VerifyPSS(&c.key.PublicKey, crypto.SHA256, h[:], s, &rsa.PSSOptions{SaltLength: 32})
}
func (c nativeTestCrypto) Encrypt(_ context.Context, _ *providers.KeyRef, _ string, _ map[string]interface{}, b []byte) ([]byte, *providers.CryptoDetails, error) {
	block, _ := aes.NewCipher(c.aes)
	g, _ := cipher.NewGCM(block)
	n := make([]byte, g.NonceSize())
	_, _ = rand.Read(n)
	return g.Seal(n, n, b, nil), nil, nil
}
func (c nativeTestCrypto) Decrypt(_ context.Context, _ *providers.KeyRef, _ string, _ map[string]interface{}, b []byte) ([]byte, error) {
	block, _ := aes.NewCipher(c.aes)
	g, _ := cipher.NewGCM(block)
	if len(b) < g.NonceSize() {
		return nil, errors.New("short")
	}
	return g.Open(nil, b[:g.NonceSize()], b[g.NonceSize():], nil)
}
func TestNativeSSOBindings(t *testing.T) {
	key, e := rsa.GenerateKey(rand.Reader, 2048)
	if e != nil {
		t.Fatal(e)
	}
	crypt := nativeTestCrypto{key: key, aes: make([]byte, 32)}
	_, _ = rand.Read(crypt.aes)
	cfg := &config.AppConfig{Issuer: "http://localhost:8088"}
	cfg.JWT.PreferredKeyID = "RSA_2048"
	now := time.Now()
	family := "vendor-a"
	parentAllowed := true
	p := &deviceSecretProvider{cfg: cfg, crypto: crypt, now: func() time.Time { return now }, sessions: map[string]int64{}, policy: func(_ context.Context, id string) (string, []string, error) {
		if id == "parent" {
			s := []string{"openid", "profile", "email"}
			if parentAllowed {
				s = append(s, "device_sso")
			}
			return "vendor-a", s, nil
		}
		return family, []string{"openid", "profile", "email"}, nil
	}}
	ctx := context.Background()
	result, e := p.GenerateDeviceSecret(ctx, providers.DeviceSecretContext{Subject: "person", ClientID: "parent", AuthID: "authentication", DPoPJkt: "device-key", Scopes: []string{"openid", "profile", "email", "device_sso"}, AuthTime: now.Unix(), ACR: "otp"})
	if e != nil {
		t.Fatal(e)
	}
	idClaims := map[string]interface{}{"iss": cfg.Issuer, "aud": "parent", "sub": "person", "sid": result.SID, "ds_hash": result.Hash, "exp": now.Unix() + 900, "iat": now.Unix()}
	sign := func(m map[string]interface{}) string {
		h := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"PS256","typ":"JWT"}`))
		b, _ := json.Marshal(m)
		in := h + "." + base64.RawURLEncoding.EncodeToString(b)
		sig, _ := crypt.Sign(ctx, providers.KeyRef{}, "PS256", []byte(in))
		return in + "." + base64.RawURLEncoding.EncodeToString(sig)
	}
	id := sign(idClaims)
	if _, e := p.ValidateDeviceSecret(ctx, result.Secret, id, "child", "device-key", []string{"openid", "email"}); e != nil {
		t.Fatal("valid exchange rejected", e)
	}
	cases := []struct {
		name, secret, id, child, jkt string
		scopes                       []string
	}{
		{"wrong device", result.Secret, id, "child", "other", []string{"openid"}},
		{"missing proof", result.Secret, id, "child", "", []string{"openid"}},
		{"scope widening", result.Secret, id, "child", "device-key", []string{"openid", "admin"}},
		{"device_sso in child", result.Secret, id, "child", "device-key", []string{"device_sso"}},
		{"empty scopes", result.Secret, id, "child", "device-key", nil},
		{"same client", result.Secret, id, "parent", "device-key", []string{"openid"}},
		{"tampered secret", "A" + result.Secret[1:], id, "child", "device-key", []string{"openid"}},
		{"tampered signature", result.Secret, id[:len(id)-10] + "AAAAAAAAAA", "child", "device-key", []string{"openid"}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, e := p.ValidateDeviceSecret(ctx, tc.secret, tc.id, tc.child, tc.jkt, tc.scopes); e == nil {
				t.Fatal("accepted invalid exchange")
			}
		})
	}
	for field, value := range map[string]interface{}{"sub": "other", "aud": "other", "sid": "other", "ds_hash": "other", "iss": "other", "exp": now.Unix() - 1} {
		t.Run("id binding "+field, func(t *testing.T) {
			old := idClaims[field]
			idClaims[field] = value
			bad := sign(idClaims)
			idClaims[field] = old
			if _, e := p.ValidateDeviceSecret(ctx, result.Secret, bad, "child", "device-key", []string{"openid"}); e == nil {
				t.Fatal("accepted mismatched " + field)
			}
		})
	}
	family = "vendor-b"
	if _, e := p.ValidateDeviceSecret(ctx, result.Secret, id, "child", "device-key", []string{"openid"}); e == nil {
		t.Fatal("cross family accepted")
	}
	family = "vendor-a"
	parentAllowed = false
	if _, e := p.ValidateDeviceSecret(ctx, result.Secret, id, "child", "device-key", []string{"openid"}); e == nil {
		t.Fatal("removed parent permission accepted")
	}
	parentAllowed = true
	p.sessions = map[string]int64{}
	if _, e := p.ValidateDeviceSecret(ctx, result.Secret, id, "child", "device-key", []string{"openid"}); e == nil {
		t.Fatal("lost session accepted")
	}
	if strings.Contains(result.Secret, "person") {
		t.Fatal("plaintext secret")
	}
}
