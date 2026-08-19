package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Middleware is an OIDC client stub adapted from pf-workspace apps/api/internal/auth.
// Dev header X-Dev-User-Sub matches P04/P06. Bearer tokens call P01 /userinfo.
// JWKS JWT parse (lestrrat-go/jwx in P04) is omitted so generated tests stay offline.

type ctxKey struct{}

type User struct {
	Sub string
}

func WithUser(ctx context.Context, u User) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

func UserFrom(ctx context.Context) (User, bool) {
	u, ok := ctx.Value(ctxKey{}).(User)
	return u, ok
}

type Middleware struct {
	devAuth      bool
	issuer       string
	internalBase string
	audience     string
}

func New(devAuth bool, issuer, internalBase, audience string) *Middleware {
	if internalBase == "" {
		internalBase = issuer
	}
	return &Middleware{devAuth: devAuth, issuer: issuer, internalBase: internalBase, audience: audience}
}

func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, err := m.authenticate(r)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithUser(r.Context(), u)))
	})
}

func (m *Middleware) authenticate(r *http.Request) (User, error) {
	if h := strings.TrimSpace(r.Header.Get("X-Dev-User-Sub")); h != "" && m.devAuth {
		return User{Sub: h}, nil
	}
	authz := strings.TrimSpace(r.Header.Get("Authorization"))
	if !strings.HasPrefix(authz, "Bearer ") {
		return User{}, fmt.Errorf("missing bearer")
	}
	token := strings.TrimSpace(strings.TrimPrefix(authz, "Bearer "))
	if token == "" {
		return User{}, fmt.Errorf("missing bearer")
	}
	if m.issuer == "" {
		return User{}, fmt.Errorf("oidc not configured")
	}
	return m.authenticateUserInfo(r.Context(), token)
}

func (m *Middleware) authenticateUserInfo(ctx context.Context, token string) (User, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.internalBase+"/userinfo", nil)
	if err != nil {
		return User{}, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return User{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("userinfo %d", res.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(res.Body, 4096))
	if err != nil {
		return User{}, err
	}
	var ui struct {
		Sub string `json:"sub"`
	}
	if err := json.Unmarshal(body, &ui); err != nil {
		return User{}, err
	}
	if ui.Sub == "" {
		return User{}, fmt.Errorf("empty sub")
	}
	_ = m.audience // reserved for JWKS/JWT verify when overlaying P01 the way P04 does
	return User{Sub: ui.Sub}, nil
}
