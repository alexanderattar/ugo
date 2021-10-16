package middleware

import (
	"net/http"

	_ "github.com/lib/pq" // postgres driver
	newrelic "github.com/newrelic/go-agent"
)

// NewRelicMiddleware adds New Relic transactions to http requests
func NewRelicMiddleware(nr newrelic.Application) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			txn := nr.StartTransaction(r.RequestURI, w, r)
			defer txn.End()
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
