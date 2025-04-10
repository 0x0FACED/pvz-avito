package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/0x0FACED/pvz-avito/internal/pkg/httpcommon"
	"github.com/0x0FACED/pvz-avito/internal/pkg/logger"
)

type Middleware struct {
	jwtManager *httpcommon.JWTManager
	log        *logger.ZerologLogger
}

func NewMiddlewareHandler(jwt *httpcommon.JWTManager, l *logger.ZerologLogger) *Middleware {
	return &Middleware{
		jwtManager: jwt,
		log:        l,
	}
}

func (m *Middleware) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			httpcommon.JSONError(w, http.StatusUnauthorized, errors.New("no auth"))
			return
		}

		parts := strings.Split(tokenString, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			httpcommon.JSONError(w, http.StatusUnauthorized, errors.New("no bearer"))
			return
		}

		claims, err := m.jwtManager.Verify(parts[1])
		if err != nil {
			httpcommon.JSONError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		m.log.Info().Str("user_email", claims.Email).Str("user_role", claims.Role).Msg("User authenticated successfully")

		ctx := context.WithValue(r.Context(), "user", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func (m *Middleware) Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		bodyBytes, _ := io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Get query params of req
		queryParams := r.URL.Query()

		// Get form data of req
		var formData map[string]any
		if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" || r.Header.Get("Content-Type") == "multipart/form-data" {
			if err := r.ParseForm(); err == nil {
				formData = make(map[string]any)
				for key, values := range r.Form {
					if len(values) == 1 {
						formData[key] = values[0]
					} else {
						formData[key] = values
					}
				}
			}
		}

		// Get JSON body of req
		var jsonBody map[string]any
		if strings.Contains(r.Header.Get("Content-Type"), "application/json") {

			// Restore req body for future using in handlers
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Parse JSON
			if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
				jsonBody = nil // just ingore it of err != nil
			}

		}

		// Serve next handler
		next.ServeHTTP(w, r)

		// Create log event
		logEvent := m.log.Info().
			Str("method", r.Method).
			Str("addr", r.RemoteAddr).
			Str("host", r.Host).
			Str("request_uri", r.RequestURI).
			TimeDiff("duration(ms)", time.Now(), start).
			Str("content_type", r.Header.Get("Content-Type"))

		// If there are query params - log it
		if len(queryParams) > 0 {
			logEvent = logEvent.Interface("query_params", queryParams)
		}

		// If there is form data - log it
		if len(formData) > 0 {
			logEvent = logEvent.Interface("form_data", formData)
		}

		// If there is json body - log it
		if len(jsonBody) > 0 {
			logEvent = logEvent.Interface("json_body", jsonBody)
		}

		// Write final log
		logEvent.Msg("Request")
	})
}
