package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Review struct {
	ID      int    `json:"id"`
	HotelID int    `json:"hotel_id"`
	User    string `json:"user"`
	Score   int    `json:"score"`
	Message string `json:"message"`
}

var seq = 2
var reviews = []Review{{ID: 1, HotelID: 1, User: "Rina", Score: 5, Message: "Hotel bersih dan nyaman"}}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/reviews/add", addReview)
	mux.HandleFunc("/api/reviews/list", listReviews)
	mux.HandleFunc("/api/reviews/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("REVIEW_SERVICE_PORT")
	if port == "" {
		port = "9107"
	}
	log.Printf("review-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func addReview(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	hotelID, _ := strconv.Atoi(r.FormValue("hotel_id"))
	score, _ := strconv.Atoi(r.FormValue("score"))
	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}
	seq++
	rv := Review{ID: seq, HotelID: hotelID, User: r.FormValue("user"), Score: score, Message: r.FormValue("message")}
	reviews = append(reviews, rv)
	writeJSON(w, map[string]any{"message": "review ditambahkan", "review": rv})
}

func listReviews(w http.ResponseWriter, r *http.Request) {
	hid, _ := strconv.Atoi(r.URL.Query().Get("hotel_id"))
	out := make([]Review, 0)
	for _, rv := range reviews {
		if hid == 0 || rv.HotelID == hid {
			out = append(out, rv)
		}
	}
	avg := 0.0
	if len(out) > 0 {
		sum := 0
		for _, x := range out {
			sum += x.Score
		}
		avg = float64(sum) / float64(len(out))
	}
	writeJSON(w, map[string]any{"count": len(out), "average": avg, "items": out})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
