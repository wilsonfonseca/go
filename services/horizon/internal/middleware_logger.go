package horizon

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/stellar/go/services/horizon/internal/render"
	"github.com/stellar/go/support/log"
)

const (
	clientNameHeader    = "X-Client-Name"
	clientVersionHeader = "X-Client-Version"
)

// LoggerMiddleware is the middleware that logs http requests and resposnes
// to the logging subsytem of horizon.
func LoggerMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		mw := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		logger := log.WithField("req", chimiddleware.GetReqID(ctx))

		ctx = log.Set(ctx, logger)

		logStartOfRequest(ctx, r)

		then := time.Now()
		h.ServeHTTP(mw, r.WithContext(ctx))
		duration := time.Now().Sub(then)
		// Checking `Accept` header from user request because if the streaming connection
		// is reset before sending the first event no Content-Type header is sent in a response.
		acceptHeader := r.Header.Get("Accept")
		streaming := strings.Contains(acceptHeader, render.MimeEventStream)
		logEndOfRequest(ctx, r, duration, mw, streaming)
	}

	return http.HandlerFunc(fn)
}

// getClientData gets client data (name or version) from header or GET parameter
// (useful when not possible to set headers, like in EventStream).
func getClientData(r *http.Request, headerName string) string {
	value := r.Header.Get(headerName)
	if value != "" {
		return value
	}

	return r.URL.Query().Get(headerName)
}

func logStartOfRequest(ctx context.Context, r *http.Request) {
	log.Ctx(ctx).WithFields(log.F{
		"client_name":    getClientData(r, clientNameHeader),
		"client_version": getClientData(r, clientVersionHeader),
		"forwarded_ip":   firstXForwardedFor(r),
		"host":           r.Host,
		"ip":             remoteAddrIP(r),
		"ip_port":        r.RemoteAddr,
		"method":         r.Method,
		"path":           r.URL.String(),
	}).Info("Starting request")
}

func logEndOfRequest(ctx context.Context, r *http.Request, duration time.Duration, mw middleware.WrapResponseWriter, streaming bool) {
	log.Ctx(ctx).WithFields(log.F{
		"bytes":          mw.BytesWritten(),
		"client_name":    getClientData(r, clientNameHeader),
		"client_version": getClientData(r, clientVersionHeader),
		"duration":       duration.Seconds(),
		"forwarded_ip":   firstXForwardedFor(r),
		"host":           r.Host,
		"ip":             remoteAddrIP(r),
		"ip_port":        r.RemoteAddr,
		"method":         r.Method,
		"path":           r.URL.String(),
		"route":          chi.RouteContext(r.Context()).RoutePattern(),
		"status":         mw.Status(),
		"streaming":      streaming,
	}).Info("Finished request")
}
