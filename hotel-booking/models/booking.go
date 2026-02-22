package models

type Booking struct {
	ID      int
	UserID  int
	RoomID  int
	Nights  int
	Total   int
	Status  string
}
