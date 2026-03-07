package controllers

import (
	"net/http"

	"hotel-booking/models"
)

type UserBookingRow struct {
	ID        int
	Status    string
	Total     int
	CheckIn   string
	CheckOut  string
	RoomID    int
	PromoCode string
}

type UserDashboardData struct {
	Title          string
	Bookings       []UserBookingRow
	WishlistHotels []models.Hotel
	ServiceOrders  []models.ServiceOrder
}

func (a *App) DashboardPage(w http.ResponseWriter, _ *http.Request) {
	bookings := make([]UserBookingRow, 0)
	for _, b := range a.DB.Bookings {
		if b.UserID != 3 {
			continue
		}
		bookings = append(bookings, UserBookingRow{ID: b.ID, Status: b.Status, Total: b.Total, CheckIn: b.CheckInDate.Format("02 Jan 2006"), CheckOut: b.CheckOutDate.Format("02 Jan 2006"), RoomID: b.RoomID, PromoCode: b.PromoCode})
	}
	wishlistHotels := make([]models.Hotel, 0)
	for _, wItem := range a.DB.Wishlists {
		if wItem.UserID != 3 {
			continue
		}
		for _, h := range a.DB.Hotels {
			if h.ID == wItem.HotelID {
				wishlistHotels = append(wishlistHotels, h)
				break
			}
		}
	}
	services := make([]models.ServiceOrder, 0)
	for _, s := range a.DB.ServiceOrders {
		services = append(services, s)
	}
	render(w, "dashboard.html", UserDashboardData{Title: "Dashboard Pengguna", Bookings: bookings, WishlistHotels: wishlistHotels, ServiceOrders: services})
}
