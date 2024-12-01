package api

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type contextKey string

const requestIDKey = contextKey("requestID")

type customResponseWriter struct {
	http.ResponseWriter
	statusCode int
	requestID  string
}

// WriteHeader Capture the status code
func (w *customResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (api *GreetingUserAPI) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				serverErrorResponse(w, r, fmt.Errorf("%s", err), api.Logger)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (api *GreetingUserAPI) logResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		ctx := context.WithValue(r.Context(), requestIDKey, requestID)
		customWriter := &customResponseWriter{ResponseWriter: w, requestID: requestID}

		start := time.Now().UTC()
		next.ServeHTTP(customWriter, r.WithContext(ctx))
		end := time.Now().UTC()
		latency := end.Sub(start)

		api.Logger.WithFields(log.Fields{
			"Latency":    latency.String(),
			"StatusCode": strconv.Itoa(customWriter.statusCode),
			"UUID":       customWriter.requestID,
		}).Info("Response Log")

	})
}

func (api *GreetingUserAPI) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		requestID := r.Context().Value(requestIDKey).(string)

		api.Logger.WithFields(log.Fields{
			"RemoteAddress":   r.RemoteAddr,
			"ProtocolVersion": r.Proto,
			"Method":          r.Method,
			"RequestURI":      r.URL.RequestURI(),
			"Host":            r.Host,
			"UUID":            requestID,
		}).Info("Request Log")

		next.ServeHTTP(w, r)
	})
}

func (api *GreetingUserAPI) SecureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy",
			"default-src 'self'; "+
				"script-src 'self' 'unsafe-inline' code.jquery.com stackpath.bootstrapcdn.com cdnjs.cloudflare.com netdna.bootstrapcdn.com; "+
				"style-src 'self' 'unsafe-inline' fonts.googleapis.com stackpath.bootstrapcdn.com cdnjs.cloudflare.com netdna.bootstrapcdn.com; "+
				"font-src 'self' fonts.gstatic.com netdna.bootstrapcdn.com;")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (api *GreetingUserAPI) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Baggage, Sentry-Trace")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
