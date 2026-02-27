package controllers

import "net/http"

func (a *App) EmailPreview(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Email reminder sent (mock)"))
}
