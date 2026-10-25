package cmodels

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/suite"
)

// fakeConfigCryptoProvider is a minimal kmprovider.ConfigCryptoProvider test double. It cannot
// use the generated cryptomock package here, since that package (transitively, via
// pkg/thunderidengine/providers) imports cmodels, which would create an import cycle.
type fakeConfigCryptoProvider struct {
	encryptFn func([]byte) ([]byte, error)
	decryptFn func([]byte) ([]byte, error)
}

func (f *fakeConfigCryptoProvider) Encrypt(_ context.Context, content []byte) ([]byte, error) {
	return f.encryptFn(content)
}

func (f *fakeConfigCryptoProvider) Decrypt(_ context.Context, content []byte) ([]byte, error) {
	return f.decryptFn(content)
}

type PropertyTestSuite struct {
	suite.Suite
}

func TestPropertyTestSuite(t *testing.T) {
	suite.Run(t, new(PropertyTestSuite))
}

func (s *PropertyTestSuite) TestSerializePropertiesToJSONObject_EmptySlice() {
	result, err := SerializePropertiesToJSONObject([]Property{})
	s.NoError(err)
	s.Equal("", result)
}

func (s *PropertyTestSuite) TestSerializePropertiesToJSONObject_NilSlice() {
	result, err := SerializePropertiesToJSONObject(nil)
	s.NoError(err)
	s.Equal("", result)
}

func (s *PropertyTestSuite) TestSerializePropertiesToJSONObject_MultipleProperties() {
	props := []Property{
		{name: "client_id", value: "my-client", isSecret: false},
		{name: "api_key", value: "secret-val", isSecret: true},
	}

	result, err := SerializePropertiesToJSONObject(props)
	s.Require().NoError(err)
	s.NotEmpty(result)

	deserialized, err := DeserializePropertiesFromJSONObject(result)
	s.Require().NoError(err)
	s.Len(deserialized, 2)

	byName := make(map[string]Property, len(deserialized))
	for _, p := range deserialized {
		byName[p.name] = p
	}

	clientProp := byName["client_id"]
	s.Equal("my-client", clientProp.value)
	s.False(clientProp.isSecret)

	apiKeyProp := byName["api_key"]
	s.Equal("secret-val", apiKeyProp.value)
	s.True(apiKeyProp.isSecret)
}

func (s *PropertyTestSuite) TestSerializePropertiesToJSONObject_PreservesIsSecretFlag() {
	props := []Property{
		{name: "secret_prop", value: "hidden", isSecret: true},
		{name: "plain_prop", value: "visible", isSecret: false},
	}

	result, err := SerializePropertiesToJSONObject(props)
	s.Require().NoError(err)

	deserialized, err := DeserializePropertiesFromJSONObject(result)
	s.Require().NoError(err)

	byName := make(map[string]Property, len(deserialized))
	for _, p := range deserialized {
		byName[p.name] = p
	}

	s.True(byName["secret_prop"].isSecret)
	s.False(byName["plain_prop"].isSecret)
}

func (s *PropertyTestSuite) TestDeserializePropertiesFromJSONObject_EmptyString() {
	result, err := DeserializePropertiesFromJSONObject("")
	s.NoError(err)
	s.Empty(result)
}

func (s *PropertyTestSuite) TestDeserializePropertiesFromJSONObject_ValidJSON() {
	jsonStr := `{"client_id":{"value":"my-client","isSecret":false},"token":{"value":"abc","isSecret":true}}`

	result, err := DeserializePropertiesFromJSONObject(jsonStr)
	s.Require().NoError(err)
	s.Len(result, 2)

	sort.Slice(result, func(i, j int) bool { return result[i].name < result[j].name })

	s.Equal("client_id", result[0].name)
	s.Equal("my-client", result[0].value)
	s.False(result[0].isSecret)

	s.Equal("token", result[1].name)
	s.Equal("abc", result[1].value)
	s.True(result[1].isSecret)
}

func (s *PropertyTestSuite) TestDeserializePropertiesFromJSONObject_InvalidJSON() {
	result, err := DeserializePropertiesFromJSONObject("{invalid")
	s.Error(err)
	s.Nil(result)
}

func (s *PropertyTestSuite) TestSerializeDeserializePropertiesFromJSONObject_Roundtrip() {
	original := []Property{
		{name: "key1", value: "val1", isSecret: false},
		{name: "key2", value: "val2", isSecret: true},
		{name: "key3", value: "val3", isSecret: false},
	}

	serialized, err := SerializePropertiesToJSONObject(original)
	s.Require().NoError(err)

	deserialized, err := DeserializePropertiesFromJSONObject(serialized)
	s.Require().NoError(err)
	s.Len(deserialized, len(original))

	sort.Slice(original, func(i, j int) bool { return original[i].name < original[j].name })
	sort.Slice(deserialized, func(i, j int) bool { return deserialized[i].name < deserialized[j].name })

	for i := range original {
		s.Equal(original[i].name, deserialized[i].name)
		s.Equal(original[i].value, deserialized[i].value)
		s.Equal(original[i].isSecret, deserialized[i].isSecret)
	}
}

func (s *PropertyTestSuite) TestGetValue_NotSecret() {
	p := &Property{name: "n", value: "plain", isSecret: false}
	v, err := p.GetValue()
	s.NoError(err)
	s.Equal("plain", v)
}

func (s *PropertyTestSuite) TestGetValue_NilConfigCryptoProvider() {
	SetConfigCryptoProvider(nil)
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "enc", isSecret: true}
	_, err := p.GetValue()
	s.EqualError(err, "config crypto provider not initialized")
}

func (s *PropertyTestSuite) TestGetValue_Success() {
	SetConfigCryptoProvider(&fakeConfigCryptoProvider{
		decryptFn: func(content []byte) ([]byte, error) { return []byte("plain"), nil },
	})
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "enc", isSecret: true}
	v, err := p.GetValue()
	s.NoError(err)
	s.Equal("plain", v)
}

func (s *PropertyTestSuite) TestGetValue_DecryptError() {
	SetConfigCryptoProvider(&fakeConfigCryptoProvider{
		decryptFn: func(content []byte) ([]byte, error) { return nil, errors.New("boom") },
	})
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "enc", isSecret: true}
	_, err := p.GetValue()
	s.Error(err)
	s.Contains(err.Error(), "failed to decrypt secret property")
}

func (s *PropertyTestSuite) TestEncrypt_NotSecretOrEmpty() {
	notSecret := &Property{name: "n", value: "v", isSecret: false}
	s.NoError(notSecret.Encrypt())

	emptySecret := &Property{name: "n", value: "", isSecret: true}
	s.NoError(emptySecret.Encrypt())
}

func (s *PropertyTestSuite) TestEncrypt_NilConfigCryptoProvider() {
	SetConfigCryptoProvider(nil)
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "plain", isSecret: true}
	s.EqualError(p.Encrypt(), "config crypto provider not initialized")
}

func (s *PropertyTestSuite) TestEncrypt_Success() {
	SetConfigCryptoProvider(&fakeConfigCryptoProvider{
		encryptFn: func(content []byte) ([]byte, error) { return []byte("enc"), nil },
	})
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "plain", isSecret: true}
	s.NoError(p.Encrypt())
	s.Equal("enc", p.value)
}

func (s *PropertyTestSuite) TestEncrypt_Error() {
	SetConfigCryptoProvider(&fakeConfigCryptoProvider{
		encryptFn: func(content []byte) ([]byte, error) { return nil, errors.New("boom") },
	})
	defer SetConfigCryptoProvider(nil)

	p := &Property{name: "n", value: "plain", isSecret: true}
	err := p.Encrypt()
	s.Error(err)
	s.Contains(err.Error(), "failed to encrypt secret property")
}

func (s *PropertyTestSuite) TestNewProperty_SecretEncryptsValue() {
	SetConfigCryptoProvider(&fakeConfigCryptoProvider{
		encryptFn: func(content []byte) ([]byte, error) { return []byte("enc"), nil },
	})
	defer SetConfigCryptoProvider(nil)

	p, err := NewProperty("n", "plain", true)
	s.Require().NoError(err)
	s.Equal("enc", p.value)
}

func (s *PropertyTestSuite) TestNewProperty_EncryptError() {
	SetConfigCryptoProvider(nil)
	defer SetConfigCryptoProvider(nil)

	_, err := NewProperty("n", "plain", true)
	s.Error(err)
	s.Contains(err.Error(), "failed to encrypt property")
}

func (s *PropertyTestSuite) TestNewProperty_NotSecret() {
	p, err := NewProperty("n", "plain", false)
	s.Require().NoError(err)
	s.Equal("plain", p.value)
}
