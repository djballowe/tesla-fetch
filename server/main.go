package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"tesla-app/server/common"
	"tesla-app/server/routes"
)

func main() {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("API_KEY")
	log.Println("Token: " + token)

	// Auth
	http.HandleFunc("/auth", getTeslaAuth.GetTeslaAuth)
	http.HandleFunc("/callback", getTeslaAuth.AuthCallBack)
	http.Handle("/.well-known/", http.StripPrefix("/.well-known/", http.FileServer(http.Dir("./.well-known"))))

	// Data
	http.HandleFunc("/car", getdata.GetCarStatus)

	fmt.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
}
