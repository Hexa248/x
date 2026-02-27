package models

import "time"

type Booking struct {
	ID       int
	UserID   int
	RoomID   int
	Nights   int
	Guests   int
	Total    int
	Status   string
	BookedAt time.Time
}
