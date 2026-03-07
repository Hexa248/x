package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type BookingStatus string

const (
	Pending    BookingStatus = "pending"
	Confirmed  BookingStatus = "confirmed"
	Paid       BookingStatus = "paid"
	CheckedIn  BookingStatus = "checked_in"
	CheckedOut BookingStatus = "checked_out"
	Cancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID           int           `json:"id"`
	UserID       int           `json:"user_id"`
	HotelID      int           `json:"hotel_id"`
	RoomID       int           `json:"room_id"`
	Nights       int           `json:"nights"`
	Guests       int           `json:"guests"`
	PromoCode    string        `json:"promo_code"`
	Discount     int           `json:"discount"`
	ServiceTotal int           `json:"service_total"`
	Total        int           `json:"total"`
	Status       BookingStatus `json:"status"`
	CheckInDate  time.Time     `json:"checkin_date"`
	CheckOutDate time.Time     `json:"checkout_date"`
	RefundAmount int           `json:"refund_amount"`
	CreatedAt    time.Time     `json:"created_at"`
}

var (
	seq       = 2
	bookings  = []Booking{{ID: 1, UserID: 3, HotelID: 1, RoomID: 10, Nights: 2, Guests: 2, Total: 960000, Status: Paid, CheckInDate: time.Now().AddDate(0, 0, 1), CheckOutDate: time.Now().AddDate(0, 0, 3), CreatedAt: time.Now()}}
	roomStock = map[int]int{10: 5, 11: 3, 12: 7}
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/bookings/create", createBooking)
	mux.HandleFunc("/api/bookings/list", listBookings)
	mux.HandleFunc("/api/bookings/cancel", cancelBooking)
	mux.HandleFunc("/api/bookings/checkin", checkIn)
	mux.HandleFunc("/api/bookings/checkout", checkOut)
	mux.HandleFunc("/api/bookings/service", addService)
	mux.HandleFunc("/api/bookings/health", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	port := os.Getenv("BOOKING_SERVICE_PORT")
	if port == "" {
		port = "9105"
	}
	log.Printf("booking-service on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	roomID := atoi(r.FormValue("room_id"), 10)
	nights := max(1, atoi(r.FormValue("nights"), 1))
	guests := max(1, atoi(r.FormValue("guests"), 1))
	base := atoi(r.FormValue("base_price"), 450000)
	checkIn := parseDate(r.FormValue("checkin"), time.Now().AddDate(0, 0, 1))
	checkOut := parseDate(r.FormValue("checkout"), checkIn.AddDate(0, 0, nights))
	promoCode := strings.ToUpper(strings.TrimSpace(r.FormValue("promo_code")))

	if roomStock[roomID] <= 0 {
		writeJSON(w, http.StatusConflict, map[string]any{"error": "stok kamar habis"})
		return
	}
	total := base * nights
	discount := calcPromo(promoCode, total)
	total -= discount
	serviceTotal := 0
	for _, svc := range r.Form["service"] {
		serviceTotal += servicePrice(svc)
	}
	total += serviceTotal

	seq++
	b := Booking{ID: seq, UserID: 3, HotelID: atoi(r.FormValue("hotel_id"), 1), RoomID: roomID, Nights: nights, Guests: guests, PromoCode: promoCode, Discount: discount, ServiceTotal: serviceTotal, Total: total, Status: Confirmed, CheckInDate: checkIn, CheckOutDate: checkOut, CreatedAt: time.Now()}
	bookings = append(bookings, b)
	roomStock[roomID]--
	writeJSON(w, http.StatusCreated, map[string]any{"message": "booking berhasil", "booking": b, "stock_left": roomStock[roomID]})
}

func listBookings(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"count": len(bookings), "items": bookings})
}

func cancelBooking(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id := atoi(r.FormValue("booking_id"), 0)
	for i := range bookings {
		if bookings[i].ID != id {
			continue
		}
		if bookings[i].Status == Cancelled {
			writeJSON(w, http.StatusConflict, map[string]any{"error": "sudah dibatalkan"})
			return
		}
		bookings[i].Status = Cancelled
		if bookings[i].CheckInDate.After(time.Now()) {
			bookings[i].RefundAmount = (bookings[i].Total * 80) / 100
		} else {
			bookings[i].RefundAmount = (bookings[i].Total * 40) / 100
		}
		roomStock[bookings[i].RoomID]++
		writeJSON(w, http.StatusOK, map[string]any{"message": "booking dibatalkan", "refund": bookings[i].RefundAmount, "stock_now": roomStock[bookings[i].RoomID]})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]any{"error": "booking tidak ditemukan"})
}

func checkIn(w http.ResponseWriter, r *http.Request) {
	statusTransition(w, r, CheckedIn)
}
func checkOut(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id := atoi(r.FormValue("booking_id"), 0)
	for i := range bookings {
		if bookings[i].ID != id {
			continue
		}
		bookings[i].Status = CheckedOut
		roomStock[bookings[i].RoomID]++
		writeJSON(w, http.StatusOK, map[string]any{"message": "check-out berhasil", "status": bookings[i].Status, "stock_now": roomStock[bookings[i].RoomID]})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]any{"error": "booking tidak ditemukan"})
}
func statusTransition(w http.ResponseWriter, r *http.Request, s BookingStatus) {
	_ = r.ParseForm()
	id := atoi(r.FormValue("booking_id"), 0)
	for i := range bookings {
		if bookings[i].ID != id {
			continue
		}
		bookings[i].Status = s
		writeJSON(w, http.StatusOK, map[string]any{"message": fmt.Sprintf("status %s", s), "status": bookings[i].Status})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]any{"error": "booking tidak ditemukan"})
}

func addService(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	id := atoi(r.FormValue("booking_id"), 0)
	svc := strings.TrimSpace(r.FormValue("service"))
	for i := range bookings {
		if bookings[i].ID != id {
			continue
		}
		add := servicePrice(svc)
		bookings[i].ServiceTotal += add
		bookings[i].Total += add
		writeJSON(w, http.StatusOK, map[string]any{"message": "layanan ditambahkan", "add": add, "total": bookings[i].Total})
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]any{"error": "booking tidak ditemukan"})
}

func calcPromo(code string, subtotal int) int {
	switch code {
	case "JALANYUK":
		return min((subtotal*8)/100, 100000)
	case "HEMAT10":
		return min((subtotal*10)/100, 150000)
	case "STAYVIP":
		return min((subtotal*12)/100, 220000)
	default:
		return 0
	}
}
func servicePrice(s string) int {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "room_service":
		return 85000
	case "laundry":
		return 60000
	case "spa":
		return 220000
	default:
		return 0
	}
}
func parseDate(raw string, def time.Time) time.Time {
	if raw == "" {
		return def
	}
	d, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return def
	}
	return d
}
func atoi(s string, d int) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return d
	}
	return v
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func writeJSON(w http.ResponseWriter, c int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(c)
	_ = json.NewEncoder(w).Encode(v)
}
