package internalhttp

import (
	"log"
	"net"
	"net/http"
	"time"

	"github.com/fixme_my_friend/hw12_13_14_15_calendar/internal/storage"
)

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(next http.Handler, logg Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		rqDuration := time.Since(start).Seconds()
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Fatalf("userip: %q is not IP:port", r.RemoteAddr)
		}

		logg.Infof("%s [%s] %s %s %s %d %f \"%s\"\n",
			ip,
			start.Format(storage.LayoutLog),
			r.Method,
			r.RequestURI,
			r.Proto,
			sw.status,
			rqDuration,
			r.UserAgent(),
		)
	})
}
