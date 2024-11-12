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
	http.HandleFunc("/auth", common.GetTeslaAuth)
	http.HandleFunc("/callback", common.AuthCallBack)

	// Data
	http.HandleFunc("/data", routes.GetCarStatus)

	// Command
	http.HandleFunc("/command", routes.IssueCommand)

	fmt.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
}
