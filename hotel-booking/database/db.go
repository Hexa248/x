package database

import (
	"fmt"
	"time"

	"hotel-booking/models"
)

type InMemoryDB struct {
	Users         []models.User
	Hotels        []models.Hotel
	Rooms         []models.Room
	Bookings      []models.Booking
	Payments      []models.Payment
	Notifications []models.Notification
}

func Seed() *InMemoryDB {
	db := &InMemoryDB{}
	db.Users = []models.User{
		{ID: 1, Name: "Super Admin", Email: "admin@hotel.com", Password: "admin123", Role: models.RoleAdmin},
		{ID: 2, Name: "Hotel Staff", Email: "staff@hotel.com", Password: "staff123", Role: models.RoleStaff},
		{ID: 3, Name: "Guest User", Email: "user@hotel.com", Password: "user123", Role: models.RoleUser},
	}

	cities := []string{"Jakarta", "Bandung", "Surabaya", "Yogyakarta", "Bali", "Lombok", "Medan", "Semarang", "Makassar", "Labuan Bajo", "Banda Aceh"}
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
		"https://images.unsplash.com/photo-1473951574080-01fe45ec8643?auto=format&fit=crop&w=1200&q=80",
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
		"https://maps.google.com/maps?q=5.548290,95.323753&z=12&output=embed",
	}

	hotelPrefixes := []string{"Nusantara", "Skyline", "Vista", "Grand", "Urban", "Harbor", "Royal", "Summit", "Lotus", "Prime", "Emerald", "Heritage", "Serene", "Aston", "Sunrise"}
	ratings := []float64{4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 4.8, 4.9}
	vipFacilities := [][]string{{"Private Jacuzzi", "Smart TV 65\"", "Lounge Access", "Mini Bar", "Bathtub"}, {"Private Sauna", "Butler Service", "Ocean View", "Premium WiFi", "Coffee Machine"}, {"Club Lounge", "King Bed", "Pillow Menu", "Free Airport Transfer", "Bathtub"}}
	deluxeFacilities := [][]string{{"City View", "King Bed", "Work Desk", "Rain Shower", "Sofa"}, {"Garden View", "Queen Bed", "Netflix TV", "Mini Fridge", "Balcony"}, {"High Floor", "Twin/Double Bed", "Smart TV", "Tea Set", "Desk Lamp"}}
	regularFacilities := [][]string{{"AC", "WiFi", "TV 42\"", "Breakfast", "Shower"}, {"AC", "WiFi", "TV", "Hot Water", "2 Bottle Water"}, {"AC", "WiFi", "Work Desk", "Toiletries", "Morning Coffee"}}

	hotelID := 1
	roomID := 1
	usedHotelNames := make(map[string]bool)
	for cityIdx, city := range cities {
		for i := 0; i < 15; i++ {
			rating := ratings[(cityIdx+i)%len(ratings)]
			hotelName := fmt.Sprintf("%s %s %d", hotelPrefixes[i%len(hotelPrefixes)], city, i+1)
			for usedHotelNames[hotelName] {
				hotelName = fmt.Sprintf("%s %s %d-%d", hotelPrefixes[i%len(hotelPrefixes)], city, i+1, hotelID)
			}
			usedHotelNames[hotelName] = true
			hotel := models.Hotel{
				ID:          hotelID,
				Name:        hotelName,
				City:        city,
				Address:     fmt.Sprintf("Jl. %s No.%d", city, 10+i),
				Description: "Hotel nyaman untuk bisnis, liburan, dan keluarga dengan fasilitas lengkap serta lokasi strategis.",
				Rating:      rating,
				MapEmbedURL: maps[cityIdx],
				Image:       cityImages[cityIdx],
				Comments: []models.Comment{
					{Author: "Rina", Message: "Lokasi strategis dan kamar bersih.", Score: 5},
					{Author: "Budi", Message: "Pelayanan ramah, check-in cepat.", Score: 4},
				},
			}
			db.Hotels = append(db.Hotels, hotel)

			for variant := 0; variant < 5; variant++ {
				base := i + 1 + (variant * 15)
				db.Rooms = append(db.Rooms,
					models.Room{ID: roomID, HotelID: hotelID, Name: fmt.Sprintf("VIP-%02d", base), Type: models.RoomVIP, PricePerNight: 2200000 + (i * 30000) + (variant * 15000), Beds: 2 + (variant % 2), Capacity: 4 + (variant % 2), Stock: 5 + (variant % 3), Facilities: vipFacilities[(cityIdx+i+variant)%len(vipFacilities)], FloorMapImage: "https://images.unsplash.com/photo-1611892440504-42a792e24d32?auto=format&fit=crop&w=1000&q=80"},
					models.Room{ID: roomID + 1, HotelID: hotelID, Name: fmt.Sprintf("DELUXE-%02d", base), Type: models.RoomDeluxe, PricePerNight: 1400000 + (i * 25000) + (variant * 12000), Beds: 1 + (variant % 3), Capacity: 2 + (variant % 4), Stock: 5 + ((variant + 1) % 3), Facilities: deluxeFacilities[(cityIdx+i+variant)%len(deluxeFacilities)], FloorMapImage: "https://images.unsplash.com/photo-1631049307264-da0ec9d70304?auto=format&fit=crop&w=1000&q=80"},
					models.Room{ID: roomID + 2, HotelID: hotelID, Name: fmt.Sprintf("REG-%02d", base), Type: models.RoomRegular, PricePerNight: 700000 + (i * 20000) + (variant * 9000), Beds: 1 + (variant % 2), Capacity: 2 + (variant % 2), Stock: 5 + ((variant + 2) % 3), Facilities: regularFacilities[(cityIdx+i+variant)%len(regularFacilities)], FloorMapImage: "https://images.unsplash.com/photo-1590490360182-c33d57733427?auto=format&fit=crop&w=1000&q=80"},
				)
				roomID += 3
			}

			hotelID++
		}
	}

	now := time.Now()
	db.Bookings = []models.Booking{
		{ID: 1, UserID: 3, RoomID: 1, Nights: 3, Guests: 2, Total: 6600000, Status: "confirmed", BookedAt: now.AddDate(0, 0, -1)},
		{ID: 2, UserID: 3, RoomID: 4, Nights: 2, Guests: 1, Total: 2900000, Status: "confirmed", BookedAt: now.AddDate(0, 0, -2)},
		{ID: 3, UserID: 3, RoomID: 7, Nights: 5, Guests: 3, Total: 10400000, Status: "paid", BookedAt: now.AddDate(0, 0, -3)},
		{ID: 4, UserID: 3, RoomID: 10, Nights: 1, Guests: 2, Total: 780000, Status: "paid", BookedAt: now.AddDate(0, 0, -5)},
		{ID: 5, UserID: 3, RoomID: 13, Nights: 4, Guests: 2, Total: 6200000, Status: "confirmed", BookedAt: now.AddDate(0, 0, -8)},
		{ID: 6, UserID: 3, RoomID: 16, Nights: 2, Guests: 1, Total: 1560000, Status: "pending", BookedAt: now.AddDate(0, 0, -10)},
	}

	db.Notifications = []models.Notification{
		{ID: 1, Role: models.RoleAdmin, Message: "Sistem notifikasi admin aktif.", CreatedAt: now.Add(-2 * time.Hour)},
		{ID: 2, Role: models.RoleStaff, Message: "Sistem notifikasi staff aktif.", CreatedAt: now.Add(-90 * time.Minute)},
	}

	for i := range db.Bookings {
		for r := range db.Rooms {
			if db.Rooms[r].ID == db.Bookings[i].RoomID && db.Rooms[r].Stock > 0 {
				db.Rooms[r].Stock--
				break
			}
		}
	}

	return db
}
