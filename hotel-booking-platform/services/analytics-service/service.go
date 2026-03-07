package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/analytics/dashboard", dashboard)
	mux.HandleFunc("/api/analytics/revenue-report", revenueReport)
	mux.HandleFunc("/api/analytics/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("ANALYTICS_SERVICE_PORT")
	if port == "" {
		port = "9110"
	}
	log.Printf("analytics-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func dashboard(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{
		"total_bookings": 1240,
		"total_revenue":  982300000,
		"occupied_rooms": 327,
		"top_hotel":      "Nusa Jakarta 1",
	})
}

func revenueReport(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]any{"monthly": []map[string]any{{"month": "2026-01", "revenue": 210000000}, {"month": "2026-02", "revenue": 250000000}, {"month": "2026-03", "revenue": 300000000}}})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
