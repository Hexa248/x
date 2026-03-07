package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "9000"
	}
	log.Printf("api-gateway running on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatal(err)
	}
}
