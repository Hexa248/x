package controllers

import "net/http"

func (a *App) PaymentPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "payment.html", map[string]any{"Title": "Payment Gateway"})
}
