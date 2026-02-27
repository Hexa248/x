package main

import (
	"log"
	"net/http"

	"hotel-booking/config"
	"hotel-booking/controllers"
	"hotel-booking/database"
	"hotel-booking/middleware"
	"hotel-booking/routes"
)

func main() {
	cfg := config.Load()
	app := &controllers.App{DB: database.Seed()}

	mux := http.NewServeMux()
	routes.AuthRoutes(mux, app)
	routes.UserRoutes(mux, app)
	routes.HotelRoutes(mux, app)
	routes.BookingRoutes(mux, app)
	routes.PaymentRoutes(mux, app)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	handler := middleware.Logger(middleware.CORS(mux))
	log.Printf("user app running on :%s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatal(err)
	}
}
