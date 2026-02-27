package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func StaffRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/staff", app.StaffDashboardPage)
	mux.HandleFunc("/staff/rooms", app.StaffRoomsPage)
}
