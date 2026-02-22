package controllers

import "hotel-booking/database"

type App struct {
	DB *database.InMemoryDB
}
