package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func BookingRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/api/book", app.CreateBooking)
	mux.HandleFunc("/api/book/cancel", app.CancelBooking)
	mux.HandleFunc("/api/book/checkin", app.StaffCheckIn)
	mux.HandleFunc("/api/book/checkout", app.StaffCheckOut)
	mux.HandleFunc("/api/review", app.AddReview)
	mux.HandleFunc("/api/wishlist/toggle", app.ToggleWishlist)
	mux.HandleFunc("/api/service/add", app.AddServiceCharge)
}
