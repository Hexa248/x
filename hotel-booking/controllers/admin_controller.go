package controllers

import (
	"net/http"
	"sort"
	"strconv"

	"hotel-booking/models"
)

type AdminRoomStockRow struct {
	RoomID    int
	HotelID   int
	HotelName string
	City      string
	RoomName  string
	RoomType  models.RoomType
	Stock     int
	Capacity  int
	Beds      int
}

type AdminDashboardData struct {
	Title            string
	Users            []models.User
	Hotels           []models.Hotel
	Rooms            []models.Room
	BookingsCount    int
	PaymentsCount    int
	Notifications    []models.Notification
	AdminNotifCount  int
	HotelsByCity     map[string][]models.Hotel
	CityOrder        []string
	RoomStockRows    []AdminRoomStockRow
	LastStockMessage string
}

func (a *App) AdminPage(w http.ResponseWriter, r *http.Request) {
	adminNotifs := notificationsByRole(a.DB.Notifications, models.RoleAdmin)
	hotelsByCity := make(map[string][]models.Hotel)
	hotelByID := make(map[int]models.Hotel)
	for _, h := range a.DB.Hotels {
		hotelsByCity[h.City] = append(hotelsByCity[h.City], h)
		hotelByID[h.ID] = h
	}
	cityOrder := make([]string, 0, len(hotelsByCity))
	for city := range hotelsByCity {
		cityOrder = append(cityOrder, city)
	}
	sort.Strings(cityOrder)

	roomRows := make([]AdminRoomStockRow, 0, len(a.DB.Rooms))
	for _, room := range a.DB.Rooms {
		h := hotelByID[room.HotelID]
		roomRows = append(roomRows, AdminRoomStockRow{
			RoomID:    room.ID,
			HotelID:   room.HotelID,
			HotelName: h.Name,
			City:      h.City,
			RoomName:  room.Name,
			RoomType:  room.Type,
			Stock:     room.Stock,
			Capacity:  room.Capacity,
			Beds:      room.Beds,
		})
	}
	sort.Slice(roomRows, func(i, j int) bool {
		if roomRows[i].City != roomRows[j].City {
			return roomRows[i].City < roomRows[j].City
		}
		if roomRows[i].HotelName != roomRows[j].HotelName {
			return roomRows[i].HotelName < roomRows[j].HotelName
		}
		return roomRows[i].RoomName < roomRows[j].RoomName
	})

	data := AdminDashboardData{
		Title:            "Admin Dashboard",
		Users:            a.DB.Users,
		Hotels:           a.DB.Hotels,
		Rooms:            a.DB.Rooms,
		BookingsCount:    len(a.DB.Bookings),
		PaymentsCount:    len(a.DB.Payments),
		Notifications:    adminNotifs,
		AdminNotifCount:  len(adminNotifs),
		HotelsByCity:     hotelsByCity,
		CityOrder:        cityOrder,
		RoomStockRows:    roomRows,
		LastStockMessage: r.URL.Query().Get("stock_message"),
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

func (a *App) UpdateRoomStock(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "request tidak valid", http.StatusBadRequest)
		return
	}
	roomID, _ := strconv.Atoi(r.FormValue("room_id"))
	stock, _ := strconv.Atoi(r.FormValue("stock"))
	if roomID <= 0 {
		http.Error(w, "room_id tidak valid", http.StatusBadRequest)
		return
	}
	if stock < 0 {
		stock = 0
	}
	ok := a.DB.SetRoomStock(roomID, stock)
	if !ok {
		http.Error(w, "kamar tidak ditemukan", http.StatusNotFound)
		return
	}
	message := "Stok kamar berhasil diperbarui"
	http.Redirect(w, r, "/admin?stock_message="+message, http.StatusFound)
}
