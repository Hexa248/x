package controllers

import (
	"net/http"
	"sort"

	"hotel-booking/models"
)

type AdminDashboardData struct {
	Title           string
	Users           []models.User
	Hotels          []models.Hotel
	Rooms           []models.Room
	BookingsCount   int
	PaymentsCount   int
	Notifications   []models.Notification
	AdminNotifCount int
	HotelsByCity    map[string][]models.Hotel
	CityOrder       []string
}

func (a *App) AdminPage(w http.ResponseWriter, _ *http.Request) {
	adminNotifs := notificationsByRole(a.DB.Notifications, models.RoleAdmin)
	hotelsByCity := make(map[string][]models.Hotel)
	for _, h := range a.DB.Hotels {
		hotelsByCity[h.City] = append(hotelsByCity[h.City], h)
	}
	cityOrder := make([]string, 0, len(hotelsByCity))
	for city := range hotelsByCity {
		cityOrder = append(cityOrder, city)
	}
	sort.Strings(cityOrder)

	data := AdminDashboardData{
		Title:           "Admin Dashboard",
		Users:           a.DB.Users,
		Hotels:          a.DB.Hotels,
		Rooms:           a.DB.Rooms,
		BookingsCount:   len(a.DB.Bookings),
		PaymentsCount:   len(a.DB.Payments),
		Notifications:   adminNotifs,
		AdminNotifCount: len(adminNotifs),
		HotelsByCity:    hotelsByCity,
		CityOrder:       cityOrder,
	}
	render(w, "admin/dashboard.html", data)
}

func (a *App) AdminUsersPage(w http.ResponseWriter, _ *http.Request) {
	adminNotifs := notificationsByRole(a.DB.Notifications, models.RoleAdmin)
	data := map[string]any{
		"Title":           "Manajemen User",
		"Users":           a.DB.Users,
		"Notifications":   adminNotifs,
		"AdminNotifCount": len(adminNotifs),
	}
	render(w, "admin/users.html", data)
}
