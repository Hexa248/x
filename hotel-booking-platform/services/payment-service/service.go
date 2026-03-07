package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Payment struct {
	BookingID int    `json:"booking_id"`
	Method    string `json:"method"`
	Status    string `json:"status"`
}

var payments []Payment

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/payments/pay", pay)
	mux.HandleFunc("/api/payments/refund", refund)
	mux.HandleFunc("/api/payments/list", list)
	mux.HandleFunc("/api/payments/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "9106"
	}
	log.Printf("payment-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func pay(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bid, _ := strconv.Atoi(r.FormValue("booking_id"))
	p := Payment{BookingID: bid, Method: r.FormValue("method"), Status: "paid"}
	payments = append(payments, p)
	writeJSON(w, map[string]any{"message": "pembayaran berhasil", "payment": p})
}

func refund(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bid, _ := strconv.Atoi(r.FormValue("booking_id"))
	p := Payment{BookingID: bid, Method: "refund", Status: "refunded"}
	payments = append(payments, p)
	writeJSON(w, map[string]any{"message": "refund diproses", "payment": p})
}

func list(w http.ResponseWriter, _ *http.Request) { writeJSON(w, map[string]any{"items": payments}) }
func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
