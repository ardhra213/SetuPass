// Copyright 2026 The ThunderID Authors
// SPDX-License-Identifier: Apache-2.0

package session

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TransportTestSuite struct {
	suite.Suite
}

func TestTransportTestSuite(t *testing.T) {
	suite.Run(t, new(TransportTestSuite))
}

func (s *TransportTestSuite) TestCookieName() {
	a := CookieName("flow-1")
	b := CookieName("flow-2")

	s.True(strings.HasPrefix(a, cookieNamePrefix))
	s.Equal(a, CookieName("flow-1"), "name must be stable for a flow")
	s.NotEqual(a, b, "different flows must get different cookie names")
}

func (s *TransportTestSuite) TestInboundHandle_HandleFor() {
	ih := InboundHandle{Cookies: map[string]string{CookieName("flow-1"): "handle-1"}}

	s.Equal("handle-1", ih.HandleFor("flow-1"))
	s.Equal("", ih.HandleFor("flow-2"))
	s.Equal("", InboundHandle{}.HandleFor("flow-1"))
}

func (s *TransportTestSuite) TestInbound_ContextRoundTrip() {
	ih := InboundHandle{Cookies: map[string]string{"a": "b"}}

	got, ok := InboundFrom(WithInbound(context.Background(), ih))
	s.True(ok)
	s.Equal(ih, got)

	_, ok = InboundFrom(context.Background())
	s.False(ok)
}

func (s *TransportTestSuite) TestCookieTransport_Read() {
	transport := NewCookieTransport(false)

	r := httptest.NewRequest(http.MethodPost, "/flow/execute", nil)
	r.AddCookie(&http.Cookie{Name: CookieName("flow-1"), Value: "handle-1"})
	r.AddCookie(&http.Cookie{Name: "unrelated", Value: "x"})

	ih := transport.Read(r)

	s.Equal("handle-1", ih.HandleFor("flow-1"))
	s.Equal("x", ih.Cookies["unrelated"])
}

func (s *TransportTestSuite) TestCookieTransport_Write() {
	transport := NewCookieTransport(true)
	w := httptest.NewRecorder()

	transport.Write(w, CookieName("flow-1"), "handle-1", time.Hour)

	cookies := w.Result().Cookies()
	s.Require().Len(cookies, 1)
	ck := cookies[0]
	s.Equal(CookieName("flow-1"), ck.Name)
	s.Equal("handle-1", ck.Value)
	s.True(ck.HttpOnly)
	s.True(ck.Secure)
	s.Equal(http.SameSiteLaxMode, ck.SameSite)
	s.Equal(3600, ck.MaxAge)
}

func (s *TransportTestSuite) TestCookieTransport_Clear() {
	transport := NewCookieTransport(false)
	w := httptest.NewRecorder()

	transport.Clear(w, CookieName("flow-1"))

	cookies := w.Result().Cookies()
	s.Require().Len(cookies, 1)
	s.Equal("", cookies[0].Value)
	s.True(cookies[0].MaxAge < 0)
}

// TestCookieTransport_RoundTrip writes a handle then reads it back through a second request,
// proving the transport's write and read agree on the cookie name.
func (s *TransportTestSuite) TestCookieTransport_RoundTrip() {
	transport := NewCookieTransport(false)
	w := httptest.NewRecorder()
	transport.Write(w, CookieName("flow-1"), "handle-xyz", time.Hour)

	r := httptest.NewRequest(http.MethodPost, "/flow/execute", nil)
	for _, ck := range w.Result().Cookies() {
		r.AddCookie(ck)
	}

	ih := transport.Read(r)
	s.Equal("handle-xyz", ih.HandleFor("flow-1"))
}
