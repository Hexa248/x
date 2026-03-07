package models

import "time"

type BookingStatus string

const (
	BookingPending    BookingStatus = "pending"
	BookingConfirmed  BookingStatus = "confirmed"
	BookingPaid       BookingStatus = "paid"
	BookingCheckedIn  BookingStatus = "checked_in"
	BookingCheckedOut BookingStatus = "checked_out"
	BookingCancelled  BookingStatus = "cancelled"
)

type Booking struct {
	ID           int
	UserID       int
	RoomID       int
	Nights       int
	Guests       int
	Total        int
	Status       string
	BookedAt     time.Time
	CheckInDate  time.Time
	CheckOutDate time.Time
	PromoCode    string
	Discount     int
	ServiceTotal int
	RefundAmount int
}
