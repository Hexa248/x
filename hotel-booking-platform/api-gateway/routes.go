package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("gateway-ok"))
	})

	mux.Handle("/api/hotels/", proxy("http://127.0.0.1:9103"))
	mux.Handle("/api/rooms/", proxy("http://127.0.0.1:9104"))
	mux.Handle("/api/bookings/", proxy("http://127.0.0.1:9105"))
	mux.Handle("/api/payments/", proxy("http://127.0.0.1:9106"))
	mux.Handle("/api/reviews/", proxy("http://127.0.0.1:9107"))
	mux.Handle("/api/promos/", proxy("http://127.0.0.1:9108"))
	mux.Handle("/api/notifications/", proxy("http://127.0.0.1:9109"))
	mux.Handle("/api/analytics/", proxy("http://127.0.0.1:9110"))
	mux.Handle("/api/auth/", proxy("http://127.0.0.1:9101"))
	mux.Handle("/api/users/", proxy("http://127.0.0.1:9102"))
}

func proxy(target string) http.Handler {
	u, _ := url.Parse(target)
	return httputil.NewSingleHostReverseProxy(u)
}
