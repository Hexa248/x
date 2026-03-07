package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
)

type WishlistItem struct{ UserID, HotelID int }

var wishlists []WishlistItem

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/users/wishlist/toggle", toggleWishlist)
	mux.HandleFunc("/api/users/wishlist/list", listWishlist)
	mux.HandleFunc("/api/users/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("USER_SERVICE_PORT")
	if port == "" {
		port = "9102"
	}
	log.Printf("user-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func toggleWishlist(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	uid, _ := strconv.Atoi(r.FormValue("user_id"))
	hid, _ := strconv.Atoi(r.FormValue("hotel_id"))
	for i := range wishlists {
		if wishlists[i].UserID == uid && wishlists[i].HotelID == hid {
			wishlists = append(wishlists[:i], wishlists[i+1:]...)
			json.NewEncoder(w).Encode(map[string]any{"wishlisted": false})
			return
		}
	}
	wishlists = append(wishlists, WishlistItem{uid, hid})
	json.NewEncoder(w).Encode(map[string]any{"wishlisted": true})
}
func listWishlist(w http.ResponseWriter, r *http.Request) {
	uid, _ := strconv.Atoi(r.URL.Query().Get("user_id"))
	out := make([]WishlistItem, 0)
	for _, x := range wishlists {
		if uid == 0 || x.UserID == uid {
			out = append(out, x)
		}
	}
	json.NewEncoder(w).Encode(map[string]any{"items": out})
}
