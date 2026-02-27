package main

import (
	"log"
	"net/http"
	"os"

	"hotel-booking/controllers"
	"hotel-booking/database"
	"hotel-booking/middleware"
	"hotel-booking/routes"
)

func main() {
	port := os.Getenv("ADMIN_PORT")
	if port == "" {
		port = "8081"
	}

	app := &controllers.App{DB: database.Seed()}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/admin", http.StatusFound)
	})
	routes.AuthRoutes(mux, app)
	routes.AdminRoutes(mux, app)
	routes.StaffRoutes(mux, app)
	routes.PaymentRoutes(mux, app)
	mux.HandleFunc("/hotels", app.HotelsPage)
	mux.HandleFunc("/hotel", app.HotelDetail)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	handler := middleware.Logger(middleware.CORS(mux))
	log.Printf("admin app running on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
