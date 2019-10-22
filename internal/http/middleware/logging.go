package middleware

import (
	"log"
	"net/http"
	"time"
)

// loggingWriter wraps a normal ResponseWriter to provide some additional
// fields for logging.
type loggingWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// Needed to actually write the request
func (w *loggingWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

// Needed to capture size of a request
func (w *loggingWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}

	n, err := w.ResponseWriter.Write(b)
	w.size += n

	return n, err
}

// Logging will log some gosh darn neat stuff.
func Logging(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := loggingWriter{ResponseWriter: w}

		next.ServeHTTP(&lw, r)

		duration := time.Now().Sub(start)
		log.Printf("%v | %v | %s%v | %v | %v | %v | %v | %v | %v\n",
			lw.status,
			r.Method,
			r.Host, r.URL,
			r.UserAgent(),
			r.RemoteAddr,
			duration,
			lw.size,
			r.Proto,
			r.Header.Get("Referer"),
		)
	}

	return http.HandlerFunc(fn)
}
