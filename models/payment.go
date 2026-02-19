package models

type Payment struct {
	ID        uint   `json:"id"`
	BookingID uint   `json:"booking_id"`
	Method    string `json:"method"`
	Status    string `json:"status"`
	Token     string `json:"token"`
}
