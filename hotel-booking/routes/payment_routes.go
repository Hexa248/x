package routes

import (
	"net/http"

	"hotel-booking/controllers"
)

func PaymentRoutes(mux *http.ServeMux, app *controllers.App) {
	mux.HandleFunc("/payment", app.PaymentPage)
	mux.HandleFunc("/webhook/payment", app.PaymentWebhook)
}
