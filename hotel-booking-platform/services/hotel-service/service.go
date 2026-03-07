package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Hotel struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	City       string   `json:"city"`
	Price      int      `json:"price"`
	Rating     float64  `json:"rating"`
	Stars      int      `json:"stars"`
	Facilities []string `json:"facilities"`
}

var hotels = []Hotel{
	{1, "Nusa Bandung 1", "Bandung", 480000, 4.6, 4, []string{"wifi", "sarapan", "parkir"}},
	{2, "Nusa Bandung 2", "Bandung", 320000, 4.2, 3, []string{"wifi", "restoran"}},
	{3, "Nusa Jakarta 1", "Jakarta", 780000, 4.8, 5, []string{"wifi", "pool", "gym", "sarapan"}},
	{4, "Nusa Jakarta 2", "Jakarta", 520000, 4.4, 4, []string{"wifi", "parkir", "restoran"}},
	{5, "Nusa Bali 1", "Bali", 900000, 4.9, 5, []string{"wifi", "pool", "spa", "sarapan"}},
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hotels/list", listHotels)
	mux.HandleFunc("/api/hotels/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("HOTEL_SERVICE_PORT")
	if port == "" {
		port = "9103"
	}
	log.Printf("hotel-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func listHotels(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("q")))
	city := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("city")))
	facility := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("facility")))
	minPrice := atoiDefault(r.URL.Query().Get("min_price"), 0)
	maxPrice := atoiDefault(r.URL.Query().Get("max_price"), 2_000_000_000)
	minRating := atofDefault(r.URL.Query().Get("min_rating"), 0)
	stars := atoiDefault(r.URL.Query().Get("stars"), 0)
	sortBy := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("sort")))

	filtered := make([]Hotel, 0)
	for _, h := range hotels {
		if q != "" && !strings.Contains(strings.ToLower(h.Name), q) && !strings.Contains(strings.ToLower(h.City), q) {
			continue
		}
		if city != "" && strings.ToLower(h.City) != city {
			continue
		}
		if h.Price < minPrice || h.Price > maxPrice {
			continue
		}
		if h.Rating < minRating {
			continue
		}
		if stars > 0 && h.Stars != stars {
			continue
		}
		if facility != "" && !contains(h.Facilities, facility) {
			continue
		}
		filtered = append(filtered, h)
	}

	sort.Slice(filtered, func(i, j int) bool {
		switch sortBy {
		case "price_asc":
			return filtered[i].Price < filtered[j].Price
		case "price_desc":
			return filtered[i].Price > filtered[j].Price
		case "rating_desc":
			return filtered[i].Rating > filtered[j].Rating
		case "rating_asc":
			return filtered[i].Rating < filtered[j].Rating
		default:
			return filtered[i].ID < filtered[j].ID
		}
	})

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"count": len(filtered), "items": filtered})
}

func contains(items []string, target string) bool {
	for _, x := range items {
		if strings.EqualFold(x, target) {
			return true
		}
	}
	return false
}

func atoiDefault(s string, d int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return v
}
func atofDefault(s string, d float64) float64 {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return d
	}
	return v
}
