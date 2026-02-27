package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"hotel-booking/models"
)

type RoomData struct {
	Title   string
	Hotel   models.Hotel
	VIP     []models.Room
	Deluxe  []models.Room
	Regular []models.Room
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

	if len(data.Deluxe) > 0 {
		base := data.Deluxe[0]
		extra := make([]models.Room, 0, 15)
		for i := 1; i <= 15; i++ {
			clone := base
			clone.Name = fmt.Sprintf("DELUXE-%02d", i+1)
			clone.PricePerNight = base.PricePerNight + (i * 35000)
			clone.Stock = 2 + (i % 5)
			clone.Beds = 1 + (i % 3)
			clone.Capacity = 2 + (i % 4)
			extra = append(extra, clone)
		}
		data.Deluxe = append(data.Deluxe, extra...)
	}

	render(w, "booking.html", data)
}
