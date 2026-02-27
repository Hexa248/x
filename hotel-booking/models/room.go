package models

type RoomType string

const (
	RoomVIP     RoomType = "VIP"
	RoomDeluxe  RoomType = "Deluxe"
	RoomRegular RoomType = "Biasa"
)

type Room struct {
	ID            int
	HotelID       int
	Name          string
	Type          RoomType
	PricePerNight int
	Beds          int
	Capacity      int
	Stock         int
	Facilities    []string
	FloorMapImage string
}
