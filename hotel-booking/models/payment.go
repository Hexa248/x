package models

type Payment struct {
	ID         int
	BookingID  int
	Method     string
	GatewayRef string
	Status     string
}
