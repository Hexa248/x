package controllers

import (
	"net/http"
	"strconv"

	"hotel-booking/models"
)

type RoomData struct {
	Title    string
	Hotel    models.Hotel
	VIP      []models.Room
	Deluxe   []models.Room
	Regular  []models.Room
}

func (a *App) BookingPage(w http.ResponseWriter, r *http.Request) {
	hotelID, _ := strconv.Atoi(r.URL.Query().Get("hotel"))
	var selected models.Hotel
	for _, h := range a.DB.Hotels {
		if h.ID == hotelID {
			selected = h
		}
	}
	data := RoomData{Title: "Pilih Kamar", Hotel: selected}
	for _, room := range a.DB.Rooms {
		if room.HotelID != hotelID {
			continue
		}
		switch room.Type {
		case models.RoomVIP:
			data.VIP = append(data.VIP, room)
		case models.RoomDeluxe:
			data.Deluxe = append(data.Deluxe, room)
		default:
			data.Regular = append(data.Regular, room)
		}
	}
	render(w, "booking.html", data)
}
