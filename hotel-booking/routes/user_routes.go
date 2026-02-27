package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func UserRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/dashboard", app.DashboardPage)
}
