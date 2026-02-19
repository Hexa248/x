package models

type Room struct {
	ID         uint    `json:"id"`
	HotelID    uint    `json:"hotel_id"`
	HotelName  string  `json:"hotel_name"`
	Type       string  `json:"type"` // VIP, Deluxe, Biasa
	Number     string  `json:"number"`
	Price      float64 `json:"price"`
	Beds       int     `json:"beds"`
	Capacity   int     `json:"capacity"`
	Facilities []string `json:"facilities"`
	MapImage   string  `json:"map_image"`
}
