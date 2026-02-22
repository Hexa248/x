package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func BookingRoutes(mux *http.ServeMux, app *controllers.App) { mux.HandleFunc("/api/book", app.CreateBooking) }
