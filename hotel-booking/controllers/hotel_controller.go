package controllers

import (
	"net/http"
	"strconv"

	"hotel-booking/models"
)

type HomeData struct {
	Title        string
	Hotels       []models.Hotel
	SelectedCity string
}

type HotelDetailData struct {
	Title string
	Hotel models.Hotel
	Rooms []models.Room
}

func (a *App) Home(w http.ResponseWriter, _ *http.Request) {
	render(w, "index.html", HomeData{Title: "Beranda", Hotels: uniqueHotelsByCity(a.DB.Hotels)})
}

func uniqueHotelsByCity(hotels []models.Hotel) []models.Hotel {
	unique := make([]models.Hotel, 0)
	seen := make(map[string]bool)
	for _, h := range hotels {
		if seen[h.City] {
			continue
		}
		seen[h.City] = true
		unique = append(unique, h)
	}
	return unique
}

func (a *App) HotelsPage(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		render(w, "hotels.html", HomeData{Title: "Daftar Hotel", Hotels: a.DB.Hotels})
		return
	}

	filtered := make([]models.Hotel, 0)
	for _, h := range a.DB.Hotels {
		if h.City == city {
			filtered = append(filtered, h)
		}
	}

	render(w, "hotels.html", HomeData{Title: "Hotel di " + city, Hotels: filtered, SelectedCity: city})
}

func (a *App) HotelDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for _, h := range a.DB.Hotels {
		if h.ID == id {
			rooms := make([]models.Room, 0)
			for _, room := range a.DB.Rooms {
				if room.HotelID == h.ID {
					rooms = append(rooms, room)
				}
			}
			render(w, "hotel_detail.html", HotelDetailData{Title: h.Name, Hotel: h, Rooms: rooms})
			return
		}
	}
	http.NotFound(w, r)
}
