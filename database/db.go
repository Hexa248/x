package database

import (
	"fmt"
	"hotel-booking/models"
)

var (
	Users    []models.User
	Hotels   []models.Hotel
	Rooms    []models.Room
	Bookings []models.Booking
	Payments []models.Payment
)

func InitData() {
	Users = []models.User{
		{ID: 1, Name: "Admin", Email: "admin@hotel.com", Password: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9", Role: "admin"}, // admin123
		{ID: 2, Name: "Staff", Email: "staff@hotel.com", Password: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9", Role: "staff"},
		{ID: 3, Name: "User", Email: "user@hotel.com", Password: "240be518fabd2724ddb6f04eeb1da5967448d7e831c08c8fa822809f74c720a9", Role: "user"},
	}

	Hotels = []models.Hotel{
		{1, "Nusantara Sky Hotel", "Jakarta", "Hotel modern di pusat kota", 4.8, -6.2000, 106.8166, []string{"Lokasi strategis", "Pelayanan ramah"}},
		{2, "Bali Sunset Resort", "Bali", "Resort view pantai", 4.9, -8.4095, 115.1889, []string{"Sunset bagus", "Kamar nyaman"}},
		{3, "Bandung Green Stay", "Bandung", "Nuansa sejuk pegunungan", 4.6, -6.9175, 107.6191, []string{"Udara segar", "Sarapan enak"}},
		{4, "Yogyakarta Heritage Inn", "Yogyakarta", "Dekat area budaya", 4.5, -7.7956, 110.3695, []string{"Dekat malioboro", "Interior estetik"}},
		{5, "Surabaya Prime Lodge", "Surabaya", "Hotel bisnis premium", 4.7, -7.2575, 112.7521, []string{"Wi-Fi cepat", "Meeting room lengkap"}},
		{6, "Lombok Ocean View", "Lombok", "View laut tropis", 4.7, -8.6500, 116.3249, []string{"Pemandangan indah", "Kolam luas"}},
		{7, "Makassar Bay Hotel", "Makassar", "Hotel dekat waterfront", 4.4, -5.1477, 119.4327, []string{"Akses mudah", "Menu seafood"}},
		{8, "Medan City Comfort", "Medan", "Nyaman untuk keluarga", 4.3, 3.5952, 98.6722, []string{"Kamar luas", "Parkir lega"}},
		{9, "Semarang Riverside", "Semarang", "View sungai dan kota", 4.5, -6.9667, 110.4167, []string{"Tempat tenang", "Dekat pusat kuliner"}},
		{10, "Malang Highland Suites", "Malang", "Suhu adem dan modern", 4.6, -7.9666, 112.6326, []string{"Instagramable", "Fasilitas lengkap"}},
	}

	var roomID uint = 1
	for _, h := range Hotels {
		for i := 1; i <= 10; i++ {
			Rooms = append(Rooms,
				models.Room{ID: roomID, HotelID: h.ID, HotelName: h.Name, Type: "VIP", Number: fmt.Sprintf("VIP-%02d", i), Price: 2200000, Beds: 2, Capacity: 4, Facilities: []string{"Private Lounge", "Bathtub", "Smart TV 65\"", "Breakfast Premium", "Mini Bar"}, MapImage: fmt.Sprintf("/assets/images/maps/hotel-%d-vip.svg", h.ID)},
			)
			roomID++
			Rooms = append(Rooms,
				models.Room{ID: roomID, HotelID: h.ID, HotelName: h.Name, Type: "Deluxe", Number: fmt.Sprintf("DLX-%02d", i), Price: 1450000, Beds: 2, Capacity: 3, Facilities: []string{"City/Pool View", "Smart TV", "Breakfast", "Shower Air Panas"}, MapImage: fmt.Sprintf("/assets/images/maps/hotel-%d-deluxe.svg", h.ID)},
			)
			roomID++
			Rooms = append(Rooms,
				models.Room{ID: roomID, HotelID: h.ID, HotelName: h.Name, Type: "Biasa", Number: fmt.Sprintf("STD-%02d", i), Price: 780000, Beds: 1, Capacity: 2, Facilities: []string{"AC", "TV", "Free Wi-Fi", "Air Mineral"}, MapImage: fmt.Sprintf("/assets/images/maps/hotel-%d-biasa.svg", h.ID)},
			)
			roomID++
		}
	}
}
