package controllers

import (
	"net/http"
	"sort"
	"time"

	"hotel-booking/models"
)

type StaffBookingRow struct {
	ID              int
	Status          string
	Guests          int
	Nights          int
	RemainingNights int
	BookedAt        string
	HotelName       string
	RoomName        string
	CheckInDate     string
	CheckOutDate    string
}

type StaffDashboardData struct {
	Title            string
	Hotels           []models.Hotel
	Rooms            []models.Room
	BookingsCount    int
	MonthlyBookings  int
	TotalGuests      int
	TotalNights      int
	ThisWeekGuests   int
	ThisWeekNights   int
	ThisWeekBookings int
	ThisMonthNights  int
	RecentBookings   []StaffBookingRow
	Notifications    []models.Notification
	StaffNotifCount  int
	HotelsByCity     map[string][]models.Hotel
	CityOrder        []string
}

func (a *App) StaffDashboardPage(w http.ResponseWriter, _ *http.Request) {
	now := time.Now()
	thisYear, thisWeek := now.ISOWeek()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	totalGuests, totalNights := 0, 0
	thisWeekGuests, thisWeekNights, thisWeekBookings := 0, 0, 0
	monthlyBookings, thisMonthNights := 0, 0

	bookings := append([]models.Booking(nil), a.DB.Bookings...)
	sort.Slice(bookings, func(i, j int) bool { return bookings[i].BookedAt.After(bookings[j].BookedAt) })

	hotelByRoom := map[int]string{}
	roomName := map[int]string{}
	for _, room := range a.DB.Rooms {
		roomName[room.ID] = room.Name
		for _, h := range a.DB.Hotels {
			if h.ID == room.HotelID {
				hotelByRoom[room.ID] = h.Name
				break
			}
		}
	}

	remainingFor := func(b models.Booking) int {
		elapsed := int(now.Truncate(24*time.Hour).Sub(b.CheckInDate.Truncate(24*time.Hour)).Hours() / 24)
		if elapsed < 0 {
			elapsed = 0
		}
		rem := b.Nights - elapsed
		if rem < 0 {
			return 0
		}
		if b.Status == string(models.BookingCheckedOut) || b.Status == string(models.BookingCancelled) {
			return 0
		}
		return rem
	}

	for _, b := range bookings {
		totalGuests += b.Guests
		totalNights += b.Nights
		y, w := b.BookedAt.ISOWeek()
		if y == thisYear && w == thisWeek {
			thisWeekBookings++
			thisWeekGuests += b.Guests
			thisWeekNights += b.Nights
		}
		if !b.BookedAt.Before(monthStart) {
			monthlyBookings++
			thisMonthNights += b.Nights
		}
	}

	recent := make([]StaffBookingRow, 0, len(bookings))
	for _, b := range bookings {
		recent = append(recent, StaffBookingRow{
			ID:              b.ID,
			Status:          b.Status,
			Guests:          b.Guests,
			Nights:          b.Nights,
			RemainingNights: remainingFor(b),
			BookedAt:        b.BookedAt.Format("02 Jan 2006"),
			HotelName:       hotelByRoom[b.RoomID],
			RoomName:        roomName[b.RoomID],
			CheckInDate:     b.CheckInDate.Format("02 Jan 2006"),
			CheckOutDate:    b.CheckOutDate.Format("02 Jan 2006"),
		})
	}

	hotelsByCity := make(map[string][]models.Hotel)
	for _, h := range a.DB.Hotels {
		hotelsByCity[h.City] = append(hotelsByCity[h.City], h)
	}
	cityOrder := make([]string, 0, len(hotelsByCity))
	for city := range hotelsByCity {
		cityOrder = append(cityOrder, city)
	}
	sort.Strings(cityOrder)

	staffNotifs := notificationsByRole(a.DB.Notifications, models.RoleStaff)
	data := StaffDashboardData{
		Title:            "Staff Hotel Dashboard",
		Hotels:           a.DB.Hotels,
		Rooms:            a.DB.Rooms,
		BookingsCount:    len(bookings),
		MonthlyBookings:  monthlyBookings,
		TotalGuests:      totalGuests,
		TotalNights:      totalNights,
		ThisWeekGuests:   thisWeekGuests,
		ThisWeekNights:   thisWeekNights,
		ThisWeekBookings: thisWeekBookings,
		ThisMonthNights:  thisMonthNights,
		RecentBookings:   recent,
		Notifications:    staffNotifs,
		StaffNotifCount:  len(staffNotifs),
		HotelsByCity:     hotelsByCity,
		CityOrder:        cityOrder,
	}
	render(w, "staff/dashboard.html", data)
}

func (a *App) StaffRoomsPage(w http.ResponseWriter, _ *http.Request) {
	staffNotifs := notificationsByRole(a.DB.Notifications, models.RoleStaff)
	data := map[string]any{
		"Title":           "Operasional Kamar",
		"Hotels":          a.DB.Hotels,
		"Rooms":           a.DB.Rooms,
		"Notifications":   staffNotifs,
		"StaffNotifCount": len(staffNotifs),
	}
	render(w, "staff/rooms.html", data)
}
