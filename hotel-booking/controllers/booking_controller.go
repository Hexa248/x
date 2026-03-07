package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
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
	promoCode := strings.TrimSpace(r.FormValue("promo_code"))
	services := r.Form["service"]
	if csv := strings.TrimSpace(r.FormValue("services")); csv != "" {
		services = append(services, strings.Split(csv, ",")...)
	}

	checkInDate, checkOutDate := parseBookingDates(r.FormValue("checkin"), r.FormValue("checkout"), nights)
	if nights <= 0 {
		nights = int(checkOutDate.Sub(checkInDate).Hours() / 24)
		if nights <= 0 {
			nights = 1
		}
	}

	room, booking, err := a.DB.CreateBooking(roomID, nights, guests)
	if err != nil {
		if err.Error() == "stok kamar habis" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	if b, _ := a.DB.FindBooking(booking.ID); b != nil {
		b.CheckInDate = checkInDate
		b.CheckOutDate = checkOutDate
		discount, applied := a.DB.ApplyPromo(promoCode, b.Total)
		b.PromoCode = applied
		b.Discount = discount
		b.Total -= discount
		if b.Total < 0 {
			b.Total = 0
		}
		if len(services) > 0 {
			a.DB.AddServices(b.ID, services)
		}
		booking = *b
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
			"message":       "Booking berhasil dibuat",
			"booking_id":    booking.ID,
			"room_id":       room.ID,
			"room_name":     room.Name,
			"stock_left":    room.Stock,
			"status":        booking.Status,
			"discount":      booking.Discount,
			"promo_code":    booking.PromoCode,
			"service_total": booking.ServiceTotal,
			"total":         booking.Total,
			"mysql_synced":  a.DB.MySQL != nil,
		})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(w, "Booking dibuat. Sisa stok kamar %s: %d", room.Name, room.Stock)
}

func (a *App) CancelBooking(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bookingID, _ := strconv.Atoi(r.FormValue("booking_id"))
	booking, err := a.DB.CancelBooking(bookingID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	msg := fmt.Sprintf("Booking #%d dibatalkan. Refund: Rp %d", booking.ID, booking.RefundAmount)
	nid := len(a.DB.Notifications) + 1
	now := time.Now()
	a.DB.Notifications = append(a.DB.Notifications,
		models.Notification{ID: nid, Role: models.RoleStaff, Message: msg, CreatedAt: now},
		models.Notification{ID: nid + 1, Role: models.RoleAdmin, Message: msg, CreatedAt: now},
	)
	respondJSON(w, map[string]any{"message": msg, "status": booking.Status, "refund": booking.RefundAmount})
}

func (a *App) StaffCheckIn(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bookingID, _ := strconv.Atoi(r.FormValue("booking_id"))
	b, _ := a.DB.FindBooking(bookingID)
	if b == nil {
		http.Error(w, "booking tidak ditemukan", http.StatusNotFound)
		return
	}
	if b.Status == string(models.BookingCancelled) || b.Status == string(models.BookingCheckedOut) {
		http.Error(w, "booking tidak bisa check-in", http.StatusConflict)
		return
	}
	b.Status = string(models.BookingCheckedIn)
	nid := len(a.DB.Notifications) + 1
	a.DB.Notifications = append(a.DB.Notifications,
		models.Notification{ID: nid, Role: models.RoleAdmin, Message: fmt.Sprintf("Booking #%d sudah check-in", bookingID), CreatedAt: time.Now()},
	)
	respondJSON(w, map[string]any{"message": "Check-in berhasil", "status": b.Status})
}

func (a *App) StaffCheckOut(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bookingID, _ := strconv.Atoi(r.FormValue("booking_id"))
	b, _ := a.DB.FindBooking(bookingID)
	if b == nil {
		http.Error(w, "booking tidak ditemukan", http.StatusNotFound)
		return
	}
	if b.Status == string(models.BookingCancelled) || b.Status == string(models.BookingCheckedOut) {
		http.Error(w, "booking tidak bisa check-out", http.StatusConflict)
		return
	}
	b.Status = string(models.BookingCheckedOut)
	if room, _ := a.DB.FindRoom(b.RoomID); room != nil {
		room.Stock++
		if a.DB.MySQL != nil {
			a.DB.MySQL.SyncRoomStock(*room)
		}
	}
	respondJSON(w, map[string]any{"message": "Check-out berhasil", "status": b.Status})
}

func (a *App) AddReview(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	hotelID, _ := strconv.Atoi(r.FormValue("hotel_id"))
	score, _ := strconv.Atoi(r.FormValue("score"))
	author := strings.TrimSpace(r.FormValue("author"))
	if author == "" {
		author = "Guest"
	}
	message := strings.TrimSpace(r.FormValue("message"))
	if message == "" {
		http.Error(w, "review kosong", http.StatusBadRequest)
		return
	}
	if ok := a.DB.AddReview(hotelID, author, message, score); !ok {
		http.Error(w, "hotel tidak ditemukan", http.StatusNotFound)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/hotel?id=%d", hotelID), http.StatusFound)
}

func (a *App) ToggleWishlist(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	hotelID, _ := strconv.Atoi(r.FormValue("hotel_id"))
	active := a.DB.ToggleWishlist(3, hotelID)
	respondJSON(w, map[string]any{"wishlisted": active, "hotel_id": hotelID})
}

func (a *App) AddServiceCharge(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	bookingID, _ := strconv.Atoi(r.FormValue("booking_id"))
	service := strings.TrimSpace(r.FormValue("service"))
	if service == "" {
		http.Error(w, "service wajib diisi", http.StatusBadRequest)
		return
	}
	total := a.DB.AddServices(bookingID, []string{service})
	respondJSON(w, map[string]any{"message": "Layanan tambahan dicatat", "added": total})
}

func parseBookingDates(checkinRaw, checkoutRaw string, nights int) (time.Time, time.Time) {
	now := time.Now()
	checkIn := now.AddDate(0, 0, 1)
	if checkinRaw != "" {
		if d, err := time.Parse("2006-01-02", checkinRaw); err == nil {
			checkIn = d
		}
	}
	checkOut := checkIn.AddDate(0, 0, nights)
	if nights <= 0 {
		checkOut = checkIn.AddDate(0, 0, 1)
	}
	if checkoutRaw != "" {
		if d, err := time.Parse("2006-01-02", checkoutRaw); err == nil {
			checkOut = d
		}
	}
	if !checkOut.After(checkIn) {
		checkOut = checkIn.AddDate(0, 0, 1)
	}
	return checkIn, checkOut
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

func respondJSON(w http.ResponseWriter, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(payload)
}
