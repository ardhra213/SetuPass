// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package cors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ParserTestSuite struct {
	suite.Suite
}

func TestParserTestSuite(t *testing.T) {
	suite.Run(t, new(ParserTestSuite))
}

func (suite *ParserTestSuite) TestValidHTTPSOrigin() {
	r, err := ParseOrigin("https://example.com")
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "https://example.com", r.Raw)
	assert.False(suite.T(), r.IsNull)
}

func (suite *ParserTestSuite) TestValidHTTPOrigin() {
	r, err := ParseOrigin("http://example.com:8080")
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "http://example.com:8080", r.Raw)
	assert.False(suite.T(), r.IsNull)
}

func (suite *ParserTestSuite) TestNullOrigin() {
	r, err := ParseOrigin("null")
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "null", r.Raw)
	assert.True(suite.T(), r.IsNull)
}

func (suite *ParserTestSuite) TestEmptyHeaderRejected() {
	_, err := ParseOrigin("")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestControlCharRejected() {
	_, err := ParseOrigin("https://example.com\r\nInjected: header")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestNULByteRejected() {
	_, err := ParseOrigin("https://example.com\x00")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestDELByteRejected() {
	_, err := ParseOrigin("https://example.com\x7f")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestUnsupportedSchemeRejected() {
	cases := []string{
		"file:///etc/passwd",
		"javascript:alert(1)",
		"ftp://example.com",
		"data:text/plain,hi",
	}
	for _, c := range cases {
		_, err := ParseOrigin(c)
		suite.Require().Error(err, c)
		assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin), c)
	}
}

func (suite *ParserTestSuite) TestMissingHostRejected() {
	_, err := ParseOrigin("https://")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestUserinfoRejected() {
	_, err := ParseOrigin("https://user:pass@example.com")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestPathRejected() {
	_, err := ParseOrigin("https://example.com/callback")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestQueryRejected() {
	_, err := ParseOrigin("https://example.com?foo=bar")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestFragmentRejected() {
	_, err := ParseOrigin("https://example.com#frag")
	suite.Require().Error(err)
	assert.True(suite.T(), errors.Is(err, ErrInvalidOrigin))
}

func (suite *ParserTestSuite) TestRawIsVerbatim() {
	// Raw must echo the input exactly so the echo path can return it byte-for-byte.
	in := "HTTPS://Example.COM:8443"
	r, err := ParseOrigin(in)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), in, r.Raw)
}

func (suite *ParserTestSuite) TestCanonicalIsPopulatedForNonNullOrigins() {
	r, err := ParseOrigin("HTTPS://Example.COM")
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "https://example.com", r.Canonical,
		"ParseOrigin should compute and stash the canonical form so the matcher does not re-parse")
}

func (suite *ParserTestSuite) TestCanonicalEmptyForNullOrigin() {
	r, err := ParseOrigin("null")
	suite.Require().NoError(err)
	assert.Empty(suite.T(), r.Canonical,
		"the null origin has no canonical form; the matcher routes it through IsNull")
}

func (suite *ParserTestSuite) TestCanonicalRoundTripsIPv6() {
	r, err := ParseOrigin("https://[::1]:8443")
	suite.Require().NoError(err)
	assert.Equal(suite.T(), "https://[::1]:8443", r.Canonical)
}
