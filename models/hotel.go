package models

type Hotel struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	City        string  `json:"city"`
	Description string  `json:"description"`
	Rating      float64 `json:"rating"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	Comments    []string `json:"comments"`
}
