package database

import (
	"fmt"
	"strings"
	"time"

	"hotel-booking/models"
)

var serviceCatalog = map[string]int{
	"room_service":   85000,
	"laundry":        60000,
	"airport_pickup": 175000,
	"spa":            220000,
}

func (db *InMemoryDB) FindPromo(code string) (models.Promo, bool) {
	for _, p := range db.Promos {
		if strings.EqualFold(p.Code, strings.TrimSpace(code)) && p.Active {
			return p, true
		}
	}
	return models.Promo{}, false
}

func (db *InMemoryDB) ApplyPromo(code string, subtotal int) (discount int, appliedCode string) {
	if strings.TrimSpace(code) == "" {
		return 0, ""
	}
	promo, ok := db.FindPromo(code)
	if !ok || subtotal < promo.MinSpend {
		return 0, ""
	}
	discount = (subtotal * promo.DiscountPct) / 100
	if discount > promo.MaxDiscount {
		discount = promo.MaxDiscount
	}
	return discount, strings.ToUpper(promo.Code)
}

func (db *InMemoryDB) AddReview(hotelID int, author, message string, score int) bool {
	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}
	for i := range db.Hotels {
		if db.Hotels[i].ID != hotelID {
			continue
		}
		db.Hotels[i].Comments = append(db.Hotels[i].Comments, models.Comment{Author: author, Message: message, Score: score})
		total := 0
		for _, c := range db.Hotels[i].Comments {
			total += c.Score
		}
		db.Hotels[i].Rating = float64(total) / float64(len(db.Hotels[i].Comments))
		return true
	}
	return false
}

func (db *InMemoryDB) ToggleWishlist(userID, hotelID int) bool {
	for i := range db.Wishlists {
		if db.Wishlists[i].UserID == userID && db.Wishlists[i].HotelID == hotelID {
			db.Wishlists = append(db.Wishlists[:i], db.Wishlists[i+1:]...)
			return false
		}
	}
	db.Wishlists = append(db.Wishlists, models.WishlistItem{UserID: userID, HotelID: hotelID})
	return true
}

func (db *InMemoryDB) AddServices(bookingID int, services []string) int {
	total := 0
	for _, s := range services {
		key := strings.TrimSpace(strings.ToLower(s))
		amount, ok := serviceCatalog[key]
		if !ok {
			continue
		}
		db.ServiceOrders = append(db.ServiceOrders, models.ServiceOrder{BookingID: bookingID, ServiceName: key, Qty: 1, Amount: amount})
		total += amount
	}
	for i := range db.Bookings {
		if db.Bookings[i].ID == bookingID {
			db.Bookings[i].ServiceTotal += total
			db.Bookings[i].Total += total
			break
		}
	}
	return total
}

func (db *InMemoryDB) FindBooking(bookingID int) (*models.Booking, int) {
	for i := range db.Bookings {
		if db.Bookings[i].ID == bookingID {
			return &db.Bookings[i], i
		}
	}
	return nil, -1
}

func (db *InMemoryDB) FindRoom(roomID int) (*models.Room, int) {
	for i := range db.Rooms {
		if db.Rooms[i].ID == roomID {
			return &db.Rooms[i], i
		}
	}
	return nil, -1
}

func (db *InMemoryDB) UpdateBookingStatus(bookingID int, status models.BookingStatus) (models.Booking, error) {
	b, _ := db.FindBooking(bookingID)
	if b == nil {
		return models.Booking{}, fmt.Errorf("booking tidak ditemukan")
	}
	b.Status = string(status)
	return *b, nil
}

func (db *InMemoryDB) CancelBooking(bookingID int) (models.Booking, error) {
	b, _ := db.FindBooking(bookingID)
	if b == nil {
		return models.Booking{}, fmt.Errorf("booking tidak ditemukan")
	}
	if b.Status == string(models.BookingCancelled) {
		return *b, fmt.Errorf("booking sudah dibatalkan")
	}
	if b.Status == string(models.BookingCheckedOut) {
		return *b, fmt.Errorf("booking sudah selesai")
	}
	b.Status = string(models.BookingCancelled)
	if b.CheckInDate.After(time.Now()) {
		b.RefundAmount = (b.Total * 80) / 100
	} else {
		b.RefundAmount = (b.Total * 40) / 100
	}
	if room, _ := db.FindRoom(b.RoomID); room != nil {
		room.Stock++
		if db.MySQL != nil {
			db.MySQL.SyncRoomStock(*room)
		}
	}
	db.Payments = append(db.Payments, models.Payment{ID: len(db.Payments) + 1, BookingID: b.ID, Method: "Refund", GatewayRef: fmt.Sprintf("RFND-%d", b.ID), Status: "refunded"})
	return *b, nil
}
