package controllers

import (
	"net/http"
	"strconv"

	"hotel-booking/models"
)

type HomeData struct {
	Title  string
	Hotels []models.Hotel
}

type HotelDetailData struct {
	Title string
	Hotel models.Hotel
}

func (a *App) Home(w http.ResponseWriter, _ *http.Request) {
	render(w, "index.html", HomeData{Title: "Beranda", Hotels: a.DB.Hotels})
}

func (a *App) HotelsPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "hotels.html", HomeData{Title: "Daftar Hotel", Hotels: a.DB.Hotels})
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
