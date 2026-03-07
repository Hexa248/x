package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/auth/login", func(w http.ResponseWriter, r *http.Request) {
		_ = r.ParseForm()
		json.NewEncoder(w).Encode(map[string]any{"token": "demo-token", "email": r.FormValue("email")})
	})
	mux.HandleFunc("/api/auth/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "9101"
	}
	log.Printf("auth-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
