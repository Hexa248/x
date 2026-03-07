package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Promo struct {
	Code        string `json:"code"`
	DiscountPct int    `json:"discount_pct"`
	MaxDiscount int    `json:"max_discount"`
	MinSpend    int    `json:"min_spend"`
}

var promos = []Promo{
	{"JALANYUK", 8, 100000, 300000},
	{"HEMAT10", 10, 150000, 500000},
	{"STAYVIP", 12, 220000, 900000},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/promos/validate", validate)
	mux.HandleFunc("/api/promos/list", list)
	mux.HandleFunc("/api/promos/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("PROMO_SERVICE_PORT")
	if port == "" {
		port = "9108"
	}
	log.Printf("promo-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func list(w http.ResponseWriter, _ *http.Request) { writeJSON(w, map[string]any{"items": promos}) }

func validate(w http.ResponseWriter, r *http.Request) {
	code := strings.ToUpper(strings.TrimSpace(r.URL.Query().Get("code")))
	subtotal, _ := strconv.Atoi(r.URL.Query().Get("subtotal"))
	for _, p := range promos {
		if p.Code != code {
			continue
		}
		if subtotal < p.MinSpend {
			writeJSON(w, map[string]any{"valid": false, "reason": "min spend belum terpenuhi"})
			return
		}
		discount := (subtotal * p.DiscountPct) / 100
		if discount > p.MaxDiscount {
			discount = p.MaxDiscount
		}
		writeJSON(w, map[string]any{"valid": true, "discount": discount, "code": p.Code})
		return
	}
	writeJSON(w, map[string]any{"valid": false, "reason": "kode tidak ditemukan"})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
