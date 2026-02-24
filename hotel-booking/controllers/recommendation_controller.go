package controllers

import "net/http"

type RecommendationData struct {
	Title  string
	Hotels []HotelStory
}

type HotelStory struct {
	ID        int
	Name      string
	City      string
	Image     string
	Rating    float64
	Tagline   string
	Story     string
	Highlight string
}

func (a *App) RecommendationPage(w http.ResponseWriter, _ *http.Request) {
	stories := make([]HotelStory, 0, len(a.DB.Hotels))
	for _, h := range a.DB.Hotels {
		stories = append(stories, HotelStory{
			ID:        h.ID,
			Name:      h.Name,
			City:      h.City,
			Image:     h.Image,
			Rating:    h.Rating,
			Tagline:   "Pilihan favorit traveler minggu ini",
			Story:     "Hotel ini dikenal karena lokasi yang dekat pusat kuliner, kamar rapi, staf cepat membantu, dan suasana yang nyaman untuk staycation atau perjalanan kerja.",
			Highlight: "Tamu paling suka: kebersihan kamar, akses transportasi, dan sarapan yang variatif.",
		})
	}
	render(w, "recommendations.html", RecommendationData{Title: "Rekomendasi Hotel", Hotels: stories})
}
