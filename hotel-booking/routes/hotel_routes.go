package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func HotelRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/", app.Home)
	mux.HandleFunc("/hotels", app.HotelsPage)
	mux.HandleFunc("/hotel", app.HotelDetail)
	mux.HandleFunc("/booking", app.BookingPage)
	mux.HandleFunc("/recommendations", app.RecommendationPage)
}
