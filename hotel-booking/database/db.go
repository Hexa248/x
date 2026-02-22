package database

import (
	"fmt"

	"hotel-booking/models"
)

type InMemoryDB struct {
	Users    []models.User
	Hotels   []models.Hotel
	Rooms    []models.Room
	Bookings []models.Booking
	Payments []models.Payment
}

func Seed() *InMemoryDB {
	db := &InMemoryDB{}
	db.Users = []models.User{
		{ID: 1, Name: "Super Admin", Email: "admin@hotel.com", Password: "admin123", Role: models.RoleAdmin},
		{ID: 2, Name: "Hotel Staff", Email: "staff@hotel.com", Password: "staff123", Role: models.RoleStaff},
		{ID: 3, Name: "Guest User", Email: "user@hotel.com", Password: "user123", Role: models.RoleUser},
	}

	cities := []string{"Jakarta", "Bandung", "Surabaya", "Yogyakarta", "Bali", "Lombok", "Medan", "Semarang", "Makassar", "Labuan Bajo"}
	cityImages := []string{
		"https://images.unsplash.com/photo-1555899434-94d1368aa7af?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1559628233-100c798642d4?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1596422846543-75c6fc197f07?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1656223216279-7dcd13767ddb?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1537996194471-e657df975ab4?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1527631746610-bca00a040d60?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1626094309830-abbb0c99da4a?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1578469645742-46cae010e5d4?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1544644181-1484b3fdfc62?auto=format&fit=crop&w=1200&q=80",
		"https://images.unsplash.com/photo-1536756958594-8fbee43ec2fc?auto=format&fit=crop&w=1200&q=80",
	}
	maps := []string{
		"https://maps.google.com/maps?q=-6.200000,106.816666&z=12&output=embed",
		"https://maps.google.com/maps?q=-6.917464,107.619123&z=12&output=embed",
		"https://maps.google.com/maps?q=-7.257472,112.752090&z=12&output=embed",
		"https://maps.google.com/maps?q=-7.795580,110.369490&z=12&output=embed",
		"https://maps.google.com/maps?q=-8.409518,115.188919&z=12&output=embed",
		"https://maps.google.com/maps?q=-8.652933,116.324944&z=12&output=embed",
		"https://maps.google.com/maps?q=3.595196,98.672226&z=12&output=embed",
		"https://maps.google.com/maps?q=-6.966667,110.416664&z=12&output=embed",
		"https://maps.google.com/maps?q=-5.147665,119.432732&z=12&output=embed",
		"https://maps.google.com/maps?q=-8.496190,119.887703&z=12&output=embed",
	}

	roomID := 1
	for i := 1; i <= 10; i++ {
		hotel := models.Hotel{
			ID:          i,
			Name:        fmt.Sprintf("Nusantara Luxe %d", i),
			City:        cities[i-1],
			Address:     fmt.Sprintf("Jl. Wisata Indah No.%d", i*7),
			Description: "Desain modern perpaduan Traveloka + Booking style, cocok untuk liburan, kerja, dan keluarga.",
			Rating:      4.2 + float64(i%3)/10,
			MapEmbedURL: maps[i-1],
			Image:       cityImages[i-1],
			Comments:    []models.Comment{{Author: "Rina", Message: "Lokasi strategis dan kamarnya bersih.", Score: 5}, {Author: "Budi", Message: "Pelayanan staf ramah sekali.", Score: 4}},
		}
		db.Hotels = append(db.Hotels, hotel)

		for t := 0; t < 10; t++ {
			db.Rooms = append(db.Rooms,
				models.Room{ID: roomID, HotelID: i, Name: fmt.Sprintf("VIP-%02d", t+1), Type: models.RoomVIP, PricePerNight: 2200000 + t*25000, Beds: 2, Capacity: 4, Facilities: []string{"Private Jacuzzi", "Smart TV 65\"", "Lounge Access", "Mini Bar", "Bathtub"}, FloorMapImage: "https://images.unsplash.com/photo-1611892440504-42a792e24d32?auto=format&fit=crop&w=1000&q=80"},
				models.Room{ID: roomID + 1, HotelID: i, Name: fmt.Sprintf("DELUXE-%02d", t+1), Type: models.RoomDeluxe, PricePerNight: 1450000 + t*20000, Beds: 2, Capacity: 3, Facilities: []string{"City View", "King Bed", "Work Desk", "Rain Shower", "Sofa"}, FloorMapImage: "https://images.unsplash.com/photo-1631049307264-da0ec9d70304?auto=format&fit=crop&w=1000&q=80"},
				models.Room{ID: roomID + 2, HotelID: i, Name: fmt.Sprintf("REG-%02d", t+1), Type: models.RoomRegular, PricePerNight: 780000 + t*15000, Beds: 1, Capacity: 2, Facilities: []string{"AC", "WiFi", "TV 42\"", "Breakfast", "Shower"}, FloorMapImage: "https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=1000&q=80"},
			)
			roomID += 3
		}
	}

	return db
}
