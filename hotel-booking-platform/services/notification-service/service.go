package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

type Notification struct {
	Channel string    `json:"channel"`
	Target  string    `json:"target"`
	Event   string    `json:"event"`
	Message string    `json:"message"`
	At      time.Time `json:"at"`
}

var items []Notification

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/notifications/send", send)
	mux.HandleFunc("/api/notifications/list", list)
	mux.HandleFunc("/api/notifications/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("NOTIFICATION_SERVICE_PORT")
	if port == "" {
		port = "9109"
	}
	log.Printf("notification-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func send(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	n := Notification{Channel: r.FormValue("channel"), Target: r.FormValue("target"), Event: r.FormValue("event"), Message: r.FormValue("message"), At: time.Now()}
	items = append(items, n)
	json.NewEncoder(w).Encode(map[string]any{"message": "notifikasi terkirim", "notification": n})
}
func list(w http.ResponseWriter, _ *http.Request) {
	_ = json.NewEncoder(w).Encode(map[string]any{"items": items})
}
