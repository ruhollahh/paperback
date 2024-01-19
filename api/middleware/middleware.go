package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ruhollahh/paperback/api/contextutil"
	"github.com/ruhollahh/paperback/api/httputil"
	"github.com/ruhollahh/paperback/internal/app/service"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
)

type Middleware struct {
	Logger         *slog.Logger
	Services       service.Services
	FormDecoder    *form.Decoder
	SessionManager *scs.SessionManager
}

func (m *Middleware) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"default-src 'self'")

		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		m.Logger.Info("received request", "ip", ip, "proto", proto, "method", method, "uri", uri)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				httputil.ServerError(m.Logger, w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := m.SessionManager.GetInt64(r.Context(), "authenticatedUserID")
		var ctx context.Context
		if id == 0 {
			ctx = contextutil.ContextSetUser(r.Context(), service.AnonymousUser)
		}

		user, err := m.Services.Users.GetByID(id)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrRecordNotFound):
				ctx = contextutil.ContextSetUser(r.Context(), service.AnonymousUser)
			default:
				httputil.ServerError(m.Logger, w, r, err)
				return
			}
		} else {
			ctx = contextutil.ContextSetUser(r.Context(), user)
		}

		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) RequireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := contextutil.ContextGetUser(r.Context())

		if service.IsAnonymous(user) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	}
}

func (m *Middleware) RequireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := contextutil.ContextGetUser(r.Context())

		if !user.Activated {
			httputil.ClientError(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return m.RequireAuthenticatedUser(fn)
}

func (m *Middleware) RequirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		user := contextutil.ContextGetUser(r.Context())

		permissions, err := m.Services.Permissions.GetAllForUser(user.ID)
		if err != nil {
			httputil.ServerError(m.Logger, w, r, err)
			return
		}

		if !service.PermissionsInclude(permissions, code) {
			httputil.ClientError(w, http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	}

	return m.RequireActivatedUser(fn)
}

func (m *Middleware) CSRF(next http.Handler) http.Handler {
	handler := nosurf.New(next)
	handler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})
	return handler
}

func (m *Middleware) CSPNonce(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var buf [16]byte

		_, err := io.ReadFull(rand.Reader, buf[:])
		if err != nil {
			panic("CSP Nonce rand.Reader failed" + err.Error())
		}

		nonce := base64.RawStdEncoding.EncodeToString(buf[:])
		ctx := contextutil.ContextSetNonce(r.Context(), nonce)
		r = r.WithContext(ctx)

		csp := []string{
			"default-src 'self'",
			fmt.Sprintf("script-src 'self' 'nonce-%s'", nonce),
			fmt.Sprintf("style-src 'self' 'nonce-%s'", nonce),
		}
		h := w.Header()
		h.Set("Content-Security-Policy", strings.Join(csp, "; "))

		next.ServeHTTP(w, r)
	}
}
