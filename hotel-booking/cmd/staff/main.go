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
	port := os.Getenv("STAFF_PORT")
	if port == "" {
		port = "8082"
	}

	app := &controllers.App{DB: database.Seed()}
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.Redirect(w, r, "/staff", http.StatusFound)
	})
	routes.AuthRoutes(mux, app)
	routes.StaffRoutes(mux, app)
	routes.AdminRoutes(mux, app)
	mux.HandleFunc("/hotel", app.HotelDetail)
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	handler := middleware.Logger(middleware.CORS(mux))
	log.Printf("staff app running on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}
