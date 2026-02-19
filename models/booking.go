package models

type Booking struct {
	ID        uint    `json:"id"`
	UserID    uint    `json:"user_id"`
	HotelID   uint    `json:"hotel_id"`
	RoomID    uint    `json:"room_id"`
	Nights    int     `json:"nights"`
	Total     float64 `json:"total"`
	Status    string  `json:"status"`
}
