package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hotel-booking/models"
)

func (a *App) CreateBooking(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "request tidak valid", http.StatusBadRequest)
		return
	}

	roomID, _ := strconv.Atoi(r.FormValue("room_id"))
	nights, _ := strconv.Atoi(r.FormValue("nights"))
	guests, _ := strconv.Atoi(r.FormValue("guests"))
	if nights <= 0 {
		nights = 1
	}
	if guests <= 0 {
		guests = 1
	}

	roomIdx := -1
	for i := range a.DB.Rooms {
		if a.DB.Rooms[i].ID == roomID {
			roomIdx = i
			break
		}
	}
	if roomIdx < 0 {
		http.Error(w, "kamar tidak ditemukan", http.StatusNotFound)
		return
	}
	if a.DB.Rooms[roomIdx].Stock <= 0 {
		http.Error(w, "stok kamar habis", http.StatusConflict)
		return
	}

	a.DB.Rooms[roomIdx].Stock--
	booking := models.Booking{
		ID:       len(a.DB.Bookings) + 1,
		UserID:   3,
		RoomID:   roomID,
		Nights:   nights,
		Guests:   guests,
		Total:    a.DB.Rooms[roomIdx].PricePerNight * nights,
		Status:   "confirmed",
		BookedAt: time.Now(),
	}
	a.DB.Bookings = append(a.DB.Bookings, booking)

	message := fmt.Sprintf("Booking baru kamar %s untuk %d hari (%d tamu). Sisa stok: %d", a.DB.Rooms[roomIdx].Name, nights, guests, a.DB.Rooms[roomIdx].Stock)
	notifID := len(a.DB.Notifications) + 1
	a.DB.Notifications = append(a.DB.Notifications,
		models.Notification{ID: notifID, Role: models.RoleStaff, Message: message, CreatedAt: time.Now()},
		models.Notification{ID: notifID + 1, Role: models.RoleAdmin, Message: message, CreatedAt: time.Now()},
	)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "Booking dibuat. Sisa stok kamar %s: %d", a.DB.Rooms[roomIdx].Name, a.DB.Rooms[roomIdx].Stock)
}
