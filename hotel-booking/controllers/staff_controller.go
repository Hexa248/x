package controllers

import "net/http"

func (a *App) StaffDashboardPage(w http.ResponseWriter, _ *http.Request) {
	data := map[string]any{
		"Title":         "Staff Hotel Dashboard",
		"Hotels":        a.DB.Hotels,
		"Rooms":         a.DB.Rooms,
		"BookingsCount": len(a.DB.Bookings),
	}
	render(w, "staff/dashboard.html", data)
}

func (a *App) StaffRoomsPage(w http.ResponseWriter, _ *http.Request) {
	data := map[string]any{
		"Title":  "Operasional Kamar",
		"Hotels": a.DB.Hotels,
		"Rooms":  a.DB.Rooms,
	}
	render(w, "staff/rooms.html", data)
}
