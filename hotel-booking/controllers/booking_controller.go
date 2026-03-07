package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"hotel-booking/database"
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

	room, booking, err := a.DB.CreateBooking(roomID, nights, guests)
	if err != nil {
		if err.Error() == "stok kamar habis" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	message := fmt.Sprintf("Booking baru kamar %s untuk %d hari (%d tamu). Sisa stok: %d", room.Name, booking.Nights, booking.Guests, room.Stock)
	notifID := len(a.DB.Notifications) + 1
	now := time.Now()
	a.DB.Notifications = append(a.DB.Notifications,
		models.Notification{ID: notifID, Role: models.RoleStaff, Message: message, CreatedAt: now},
		models.Notification{ID: notifID + 1, Role: models.RoleAdmin, Message: message, CreatedAt: now},
	)

	if acceptsJSON(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"message":      "Booking berhasil dibuat",
			"booking_id":   booking.ID,
			"room_id":      room.ID,
			"room_name":    room.Name,
			"stock_left":   room.Stock,
			"status":       booking.Status,
			"mysql_synced": a.DB.MySQL != nil,
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "Booking dibuat. Sisa stok kamar %s: %d", room.Name, room.Stock)
}

func acceptsJSON(r *http.Request) bool {
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	if r.Header.Get("Accept") == "application/json" {
		return true
	}
	return database.ParseIntSafe(r.FormValue("ajax"), 0) == 1
}
