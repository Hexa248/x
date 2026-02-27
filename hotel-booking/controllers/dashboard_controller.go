package controllers

import "net/http"

func (a *App) DashboardPage(w http.ResponseWriter, _ *http.Request) {
	render(w, "dashboard.html", map[string]any{"Title": "Dashboard"})
}
