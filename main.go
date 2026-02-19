package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"hotel-booking/config"
	"hotel-booking/database"
	"hotel-booking/helpers"
	"hotel-booking/models"
)

type HomeData struct {
	Hotels          []models.Hotel
	CityCards       []CityCard
	Recommendations []models.Hotel
	SelectedCity    string
	Query           string
}

type CityCard struct {
	Name      string
	StartFrom int
	ImageURL  string
}

func main() {
	cfg := config.Load()
	database.InitData()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/hotels", hotels)
	http.HandleFunc("/hotels/", hotelActions)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/dashboard", dashboard)
	http.HandleFunc("/auth/google", googleLogin)
	http.HandleFunc("/booking", createBooking)
	http.HandleFunc("/payment/", paymentPage)
	http.HandleFunc("/payment", processPayment)

	fmt.Println("Running on http://localhost:" + cfg.Port)
	_ = http.ListenAndServe(":"+cfg.Port, nil)
}

func render(w http.ResponseWriter, tmpl string, data any) {
	_ = template.Must(template.ParseFiles("views/"+tmpl)).Execute(w, data)
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	hotels := filterHotels(query)
	render(w, "index.html", HomeData{
		Hotels:          hotels,
		CityCards:       buildCityCards(),
		Recommendations: topRecommendations(6),
		SelectedCity:    query,
		Query:           query,
	})
}

func hotels(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hotels" {
		http.NotFound(w, r)
		return
	}
	query := strings.TrimSpace(r.URL.Query().Get("city"))
	render(w, "hotels.html", map[string]any{"hotels": filterHotels(query), "city": query})
}

func hotelActions(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) < 2 {
		http.NotFound(w, r)
		return
	}
	id, _ := strconv.Atoi(parts[1])
	city := r.URL.Query().Get("city")
	var hotel models.Hotel
	for _, h := range database.Hotels {
		if int(h.ID) == id {
			hotel = h
			break
		}
	}
	if hotel.ID == 0 {
		http.NotFound(w, r)
		return
	}
	if len(parts) == 2 {
		render(w, "hotel_detail.html", map[string]any{"hotel": hotel, "city": city})
		return
	}
	if len(parts) == 3 && parts[2] == "rooms" {
		typeRoom := r.URL.Query().Get("type")
		if typeRoom == "" {
			typeRoom = "VIP"
		}
		var rooms []models.Room
		for _, rm := range database.Rooms {
			if int(rm.HotelID) == id && rm.Type == typeRoom {
				rooms = append(rooms, rm)
			}
		}
		render(w, "booking.html", map[string]any{"rooms": rooms, "roomType": typeRoom, "hotel": hotel, "city": city})
		return
	}
	http.NotFound(w, r)
}

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "login.html", nil)
		return
	}
	email, pass := r.FormValue("email"), r.FormValue("password")
	for _, u := range database.Users {
		if u.Email == email && helpers.CheckPasswordHash(pass, u.Password) {
			http.SetCookie(w, &http.Cookie{Name: "email", Value: u.Email, Path: "/"})
			http.SetCookie(w, &http.Cookie{Name: "role", Value: u.Role, Path: "/"})
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
	}
	render(w, "login.html", map[string]any{"error": "email/password salah"})
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "register.html", nil)
		return
	}
	h, _ := helpers.HashPassword(r.FormValue("password"))
	database.Users = append(database.Users, models.User{ID: uint(len(database.Users) + 1), Name: r.FormValue("name"), Email: r.FormValue("email"), Password: h, Role: "user"})
	http.Redirect(w, r, "/login", http.StatusFound)
}

func logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "email", Value: "", Path: "/", MaxAge: -1})
	http.Redirect(w, r, "/", http.StatusFound)
}

func dashboard(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("email"); err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}
	render(w, "dashboard.html", map[string]any{"bookings": database.Bookings})
}

func createBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	roomID, _ := strconv.Atoi(r.FormValue("room_id"))
	nights, _ := strconv.Atoi(r.FormValue("nights"))
	if nights < 1 {
		nights = 1
	}
	for _, rm := range database.Rooms {
		if int(rm.ID) == roomID {
			b := models.Booking{ID: uint(len(database.Bookings) + 1), UserID: 3, HotelID: rm.HotelID, RoomID: rm.ID, Nights: nights, Total: rm.Price * float64(nights), Status: "pending"}
			database.Bookings = append(database.Bookings, b)
			http.Redirect(w, r, "/payment/"+strconv.Itoa(int(b.ID)), http.StatusFound)
			return
		}
	}
	http.Error(w, "room not found", 400)
}

func paymentPage(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/payment/"))
	for _, b := range database.Bookings {
		if int(b.ID) == id {
			render(w, "payment.html", map[string]any{"booking": b})
			return
		}
	}
	http.NotFound(w, r)
}

func processPayment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method", 405)
		return
	}
	bid, _ := strconv.Atoi(r.FormValue("booking_id"))
	pay := models.Payment{ID: uint(len(database.Payments) + 1), BookingID: uint(bid), Method: r.FormValue("method"), Status: "paid", Token: helpers.CreatePaymentToken(uint(bid))}
	database.Payments = append(database.Payments, pay)
	for i := range database.Bookings {
		if int(database.Bookings[i].ID) == bid {
			database.Bookings[i].Status = "paid"
		}
	}
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

func googleLogin(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Google OAuth siap dipakai jika kredensial resmi di-set di .env"))
}

func filterHotels(query string) []models.Hotel {
	if query == "" {
		return database.Hotels
	}
	q := strings.ToLower(query)
	var out []models.Hotel
	for _, h := range database.Hotels {
		if strings.Contains(strings.ToLower(h.City), q) || strings.Contains(strings.ToLower(h.Name), q) {
			out = append(out, h)
		}
	}
	return out
}

func topRecommendations(limit int) []models.Hotel {
	h := append([]models.Hotel{}, database.Hotels...)
	sort.Slice(h, func(i, j int) bool { return h[i].Rating > h[j].Rating })
	if len(h) > limit {
		return h[:limit]
	}
	return h
}

func buildCityCards() []CityCard {
	images := []string{
		"https://images.unsplash.com/photo-1566073771259-6a8506099945?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1499856871958-5b9627545d1a?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1519817650390-64a93db51149?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1505764706515-aa95265c5abc?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1540541338287-41700207dee6?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1476514525535-07fb3b4ae5f1?auto=format&fit=crop&w=1200&q=80",
	}
	seen := map[string]bool{}
	var cards []CityCard
	idx := 0
	for _, h := range database.Hotels {
		if seen[h.City] {
			continue
		}
		seen[h.City] = true
		cards = append(cards, CityCard{Name: h.City, StartFrom: 400000 + idx*70000, ImageURL: images[idx%len(images)]})
		idx++
		if len(cards) == 6 {
			break
		}
	}
	return cards
}
