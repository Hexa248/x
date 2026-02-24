package controllers

import "net/http"

func (a *App) UserProfile(w http.ResponseWriter, _ *http.Request) {
	render(w, "dashboard.html", map[string]any{"Title": "Profil Pengguna"})
}
