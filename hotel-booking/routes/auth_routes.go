package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func AuthRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/login", app.LoginPage)
	mux.HandleFunc("/register", app.RegisterPage)
}
