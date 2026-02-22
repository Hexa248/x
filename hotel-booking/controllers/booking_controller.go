package controllers

import "net/http"

func (a *App) CreateBooking(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte("Booking berhasil dibuat (mock)."))
}
