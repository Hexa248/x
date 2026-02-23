package controllers

import (
	"fmt"
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
}

func (a *App) Home(w http.ResponseWriter, _ *http.Request) {
	render(w, "index.html", HomeData{Title: "Beranda", Hotels: a.DB.Hotels})
}

func (a *App) HotelsPage(w http.ResponseWriter, r *http.Request) {
	city := r.URL.Query().Get("city")
	if city == "" {
		render(w, "hotels.html", HomeData{Title: "Daftar Hotel", Hotels: a.DB.Hotels})
		return
	}

	for _, h := range a.DB.Hotels {
		if h.City != city {
			continue
		}
		generated := make([]models.Hotel, 0, 5)
		for i := 1; i <= 5; i++ {
			clone := h
			clone.Name = fmt.Sprintf("%s %d", h.Name, i)
			clone.Address = fmt.Sprintf("%s Tower %d", h.Address, i)
			generated = append(generated, clone)
		}
		render(w, "hotels.html", HomeData{Title: "Hotel di " + city, Hotels: generated, SelectedCity: city})
		return
	}

	render(w, "hotels.html", HomeData{Title: "Hotel", Hotels: []models.Hotel{}, SelectedCity: city})
}

func (a *App) HotelDetail(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	for _, h := range a.DB.Hotels {
		if h.ID == id {
			render(w, "hotel_detail.html", HotelDetailData{Title: h.Name, Hotel: h})
			return
		}
	}
	http.NotFound(w, r)
}
