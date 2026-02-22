package controllers

import "net/http"

func (a *App) PaymentWebhook(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("webhook received"))
}
