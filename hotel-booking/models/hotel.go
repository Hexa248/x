package models

type Hotel struct {
	ID          int
	Name        string
	City        string
	Address     string
	Description string
	Rating      float64
	MapEmbedURL string
	Image       string
	Comments    []Comment
}

type Comment struct {
	Author  string
	Message string
	Score   int
}
