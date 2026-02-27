package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func AdminRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/admin", app.AdminPage)
	mux.HandleFunc("/admin/users", app.AdminUsersPage)
}
