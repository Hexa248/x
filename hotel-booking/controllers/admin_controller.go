package controllers

import "net/http"

func (a *App) AdminPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "dashboard.html", map[string]any{"Title": "Admin Panel"})
}
