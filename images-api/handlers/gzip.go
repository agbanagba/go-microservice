package handlers

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// GzipHandler ...
type GzipHandler struct{}

// GzipResponseWriter wraps the http ResponseWriter with a gzip writer
type GzipResponseWriter struct {
	rw         http.ResponseWriter
	gzipWriter gzip.Writer
}

// NewGzipResponseWriter creates a new gzip response writer
func NewGzipResponseWriter(rw http.ResponseWriter) *GzipResponseWriter {
	gw := gzip.NewWriter(rw)
	return &GzipResponseWriter{rw, *gw}
}

// Header returns the header info of an http response
func (gziprw *GzipResponseWriter) Header() http.Header {
	return gziprw.rw.Header()
}

func (gziprw *GzipResponseWriter) Write(d []byte) (int, error) {
	return gziprw.gzipWriter.Write(d)
}

// WriteHeader of the GzipResponseWriter uses the http ResponseWriter header to write header
// information
func (gziprw *GzipResponseWriter) WriteHeader(statuscode int) {
	gziprw.rw.WriteHeader(statuscode)
}

// Flush flushes and closes the gzip writer
func (gziprw *GzipResponseWriter) Flush() {
	gziprw.gzipWriter.Flush()
	gziprw.gzipWriter.Close()
}

// GzipMiddleware ...
func (g *GzipHandler) GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			gziprw := NewGzipResponseWriter(rw)
			gziprw.Header().Set("Content-Encoding", "gzip")

			next.ServeHTTP(gziprw, r)
			defer gziprw.Flush()
		}
	})
}
