package server

import "net/http"

// contentSecurityPolicy only allows the SPA's own files. Inline styles are allowed because
// `style="..."` attributes in Svelte markup are parsed as HTML, which CSP would block; inline
// scripts are not. img-src allows data: URIs, which Vite inlines for small images.
const contentSecurityPolicy = "default-src 'self'; " +
	"script-src 'self'; " +
	"style-src 'self' 'unsafe-inline'; " +
	"img-src 'self' data:; " +
	"object-src 'none'; " +
	"base-uri 'self'; " +
	"form-action 'self'; " +
	"frame-ancestors 'none'"

// securityHeaders sets browser security headers on every response. HSTS is left to the TLS
// terminator (reverse proxy), since the app itself serves plain HTTP.
func securityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", contentSecurityPolicy)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY") // for browsers without frame-ancestors
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}
