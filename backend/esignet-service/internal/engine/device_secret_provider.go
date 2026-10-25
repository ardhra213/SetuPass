// SPDX-License-Identifier: MPL-2.0
package engine

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mosip/esignet/internal/clientmgmt"
	"github.com/mosip/esignet/internal/config"
	"github.com/mosip/esignet/internal/keymanager/signature"
	"github.com/thunder-id/thunderid/pkg/thunderidengine/providers"
)

var errNativeSSO = errors.New("native SSO validation failed")

type deviceSecretProvider struct {
	cfg      *config.AppConfig
	crypto   providers.RuntimeCryptoProvider
	policy   func(context.Context, string) (string, []string, error)
	mu       sync.Mutex
	sessions map[string]int64
	now      func() time.Time
}

// NewDeviceSecretProvider uses existing persisted keystore crypto. Family enrollment is operator-only:
// native_sso_family is intentionally NOT accepted by the public client registration/update API.
func NewDeviceSecretProvider(cfg *config.AppConfig, clients *clientmgmt.Service, crypto providers.RuntimeCryptoProvider) providers.DeviceSecretProvider {
	p := &deviceSecretProvider{cfg: cfg, crypto: crypto, sessions: map[string]int64{}, now: time.Now}
	p.policy = func(ctx context.Context, id string) (string, []string, error) {
		c, e := clients.GetActiveClient(ctx, id)
		if e != nil {
			return "", nil, e
		}
		family, _ := c.AdditionalConfig["native_sso_family"].(string)
		return family, getAllowedScopes(cfg.ScopeClaims, c.AdditionalConfig), nil
	}
	return p
}
func secretHash(s string) string {
	h := sha256.Sum256([]byte(s))
	return base64.RawURLEncoding.EncodeToString(h[:])
}
func (p *deviceSecretProvider) alg() string {
	return signature.AlgorithmForRefID(p.cfg.JWT.PreferredKeyID)
}
func (p *deviceSecretProvider) GenerateDeviceSecret(ctx context.Context, c providers.DeviceSecretContext) (*providers.DeviceSecretResult, error) {
	family, allowed, e := p.policy(ctx, c.ClientID)
	if e != nil || family == "" || !slices.Contains(allowed, "device_sso") || !slices.Contains(c.Scopes, "device_sso") || c.DPoPJkt == "" || c.Subject == "" || c.AuthID == "" {
		return nil, errNativeSSO
	}
	c.SID = uuid.NewString()
	c.ExpiresAt = p.now().Unix() + 900
	claims, _ := json.Marshal(c)
	var m map[string]interface{}
	_ = json.Unmarshal(claims, &m)
	m["iss"] = p.cfg.Issuer
	m["iat"] = p.now().Unix()
	m["cnf"] = map[string]string{"jkt": c.DPoPJkt}
	payload, e := json.Marshal(m)
	if e != nil {
		return nil, e
	}
	header, _ := json.Marshal(map[string]string{"alg": p.alg(), "typ": "device-secret+jwt"})
	enc := base64.RawURLEncoding
	input := enc.EncodeToString(header) + "." + enc.EncodeToString(payload)
	sig, e := p.crypto.Sign(ctx, providers.KeyRef{}, p.alg(), []byte(input))
	if e != nil {
		return nil, e
	}
	encrypted, _, e := p.crypto.Encrypt(ctx, nil, "AES-GCM", nil, []byte(input+"."+enc.EncodeToString(sig)))
	if e != nil {
		return nil, e
	}
	secret := enc.EncodeToString(encrypted)
	p.mu.Lock()
	defer p.mu.Unlock()
	for sid, exp := range p.sessions {
		if exp <= p.now().Unix() {
			delete(p.sessions, sid)
		}
	}
	if len(p.sessions) >= 10000 {
		return nil, errors.New("native SSO session capacity reached")
	}
	p.sessions[c.SID] = c.ExpiresAt
	return &providers.DeviceSecretResult{Secret: secret, Hash: secretHash(secret), SID: c.SID}, nil
}

// verifyJWT pins the expected server signing algorithm/key. It never trusts header-provided keys.
func (p *deviceSecretProvider) verifyJWT(ctx context.Context, token string, typ string) (map[string]interface{}, error) {
	if len(token) > 131072 {
		return nil, errNativeSSO
	}
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errNativeSSO
	}
	dec := base64.RawURLEncoding.DecodeString
	h, e := dec(parts[0])
	if e != nil {
		return nil, e
	}
	var header map[string]interface{}
	if json.Unmarshal(h, &header) != nil || header["alg"] != p.alg() || header["crit"] != nil || header["typ"] != typ {
		return nil, errNativeSSO
	}
	sig, e := dec(parts[2])
	if e != nil {
		return nil, e
	}
	if e = p.crypto.Verify(ctx, providers.KeyRef{}, p.alg(), []byte(parts[0]+"."+parts[1]), sig); e != nil {
		return nil, e
	}
	b, e := dec(parts[1])
	if e != nil {
		return nil, e
	}
	var c map[string]interface{}
	if json.Unmarshal(b, &c) != nil {
		return nil, errNativeSSO
	}
	exp, _ := c["exp"].(float64)
	iat, _ := c["iat"].(float64)
	nbf, _ := c["nbf"].(float64)
	if c["iss"] != p.cfg.Issuer || exp <= float64(p.now().Unix()) || iat > float64(p.now().Unix()+30) || nbf > float64(p.now().Unix()+30) {
		return nil, errNativeSSO
	}
	return c, nil
}
func (p *deviceSecretProvider) ValidateDeviceSecret(ctx context.Context, secret, idToken, child, jkt string, scopes []string) (*providers.DeviceSecretContext, error) {
	if len(secret) > 262144 || jkt == "" || len(scopes) == 0 || slices.Contains(scopes, "device_sso") {
		return nil, errNativeSSO
	}
	cipher, e := base64.RawURLEncoding.DecodeString(secret)
	if e != nil {
		return nil, errNativeSSO
	}
	plain, e := p.crypto.Decrypt(ctx, nil, "AES-GCM", nil, cipher)
	if e != nil {
		return nil, errNativeSSO
	}
	claims, e := p.verifyJWT(ctx, string(plain), "device-secret+jwt")
	if e != nil {
		return nil, errNativeSSO
	}
	raw, _ := json.Marshal(claims)
	var c providers.DeviceSecretContext
	if json.Unmarshal(raw, &c) != nil || c.SID == "" || c.AuthID == "" || c.Subject == "" || c.ClientID == child || c.DPoPJkt != jkt {
		return nil, errNativeSSO
	}
	cnf, ok := claims["cnf"].(map[string]interface{})
	if !ok || cnf["jkt"] != jkt {
		return nil, errNativeSSO
	}
	p.mu.Lock()
	exp, live := p.sessions[c.SID]
	p.mu.Unlock()
	if !live || exp != c.ExpiresAt || exp <= p.now().Unix() {
		return nil, errNativeSSO
	}
	parentFamily, parentScopes, e := p.policy(ctx, c.ClientID)
	if e != nil {
		return nil, errNativeSSO
	}
	childFamily, childScopes, e := p.policy(ctx, child)
	if e != nil || parentFamily == "" || parentFamily != childFamily || !slices.Contains(parentScopes, "device_sso") {
		return nil, errNativeSSO
	}
	for _, scope := range scopes {
		if !slices.Contains(c.Scopes, scope) || !slices.Contains(parentScopes, scope) || !slices.Contains(childScopes, scope) {
			return nil, errNativeSSO
		}
	}
	id, e := p.verifyJWT(ctx, idToken, "JWT")
	if e != nil {
		return nil, errNativeSSO
	}
	audOK := id["aud"] == c.ClientID
	if auds, ok := id["aud"].([]interface{}); ok {
		audOK = len(auds) == 1 && auds[0] == c.ClientID
	}
	hash, _ := id["ds_hash"].(string)
	if !audOK || id["sub"] != c.Subject || id["sid"] != c.SID || subtle.ConstantTimeCompare([]byte(hash), []byte(secretHash(secret))) != 1 {
		return nil, errNativeSSO
	}
	if azp, ok := id["azp"]; ok && azp != c.ClientID {
		return nil, errNativeSSO
	}
	c.ExpiresAt = min(c.ExpiresAt, int64(id["exp"].(float64)))
	return &c, nil
}
