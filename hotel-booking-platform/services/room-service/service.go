package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Room struct {
	ID      int    `json:"id"`
	HotelID int    `json:"hotel_id"`
	Type    string `json:"type"`
	Stock   int    `json:"stock"`
}

var rooms = []Room{{10, 1, "Deluxe", 5}, {11, 1, "VIP", 3}, {12, 2, "Regular", 7}}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/rooms/list", list)
	mux.HandleFunc("/api/rooms/stock/update", updateStock)
	mux.HandleFunc("/api/rooms/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("ROOM_SERVICE_PORT")
	if port == "" {
		port = "9104"
	}
	log.Printf("room-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func list(w http.ResponseWriter, _ *http.Request) { writeJSON(w, map[string]any{"items": rooms}) }
func updateStock(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	rid, _ := strconv.Atoi(r.FormValue("room_id"))
	s, _ := strconv.Atoi(r.FormValue("stock"))
	for i := range rooms {
		if rooms[i].ID == rid {
			rooms[i].Stock = s
			writeJSON(w, map[string]any{"message": "stok diperbarui", "room": rooms[i]})
			return
		}
	}
	writeJSON(w, map[string]any{"error": "room tidak ditemukan"})
}
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
