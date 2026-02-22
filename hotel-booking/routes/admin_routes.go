package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func AdminRoutes(mux *http.ServeMux, app *controllers.App) { mux.HandleFunc("/admin", app.AdminPage) }
