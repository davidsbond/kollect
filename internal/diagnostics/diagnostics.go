// Package diagnostics provides functions for enabling diagnostics on the application. Primarily this is for
// serving pprof endpoints.
package diagnostics

import (
	"net/http"
	"net/http/pprof"
)

// Serve diagnostic endpoints on the provided http.ServeMux.
func Serve(r *http.ServeMux) {
	if !enableDiagnostics {
		return
	}

	r.HandleFunc("/__/pprof/profile", pprof.Profile)
	r.HandleFunc("/__/pprof/trace", pprof.Trace)
	r.HandleFunc("/__/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/__/pprof/symbol", pprof.Symbol)
}
