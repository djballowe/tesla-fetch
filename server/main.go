package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"tesla-app/server/common"
	"tesla-app/server/routes"
)

func main() {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	// Auth
<<<<<<< Updated upstream
	http.HandleFunc("/auth", getTeslaAuth.GetTeslaAuth)
	http.HandleFunc("/callback", getTeslaAuth.AuthCallBack)
	http.Handle("/.well-known/", http.StripPrefix("/.well-known/", http.FileServer(http.Dir("./.well-known"))))
=======
	http.HandleFunc("/auth", common.GetTeslaAuth)
	http.HandleFunc("/callback", common.AuthCallBack)
	http.HandleFunc("/honk", getdata.GetChargeState)
>>>>>>> Stashed changes

	// Data
	http.HandleFunc("/car", getdata.GetCarStatus)

	fmt.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
}
