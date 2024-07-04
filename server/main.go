package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"tesla-app/server/routes"
	"github.com/joho/godotenv"
)

func main() {

	envErr := godotenv.Load(".env")
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("API_KEY")
	log.Println("Token: " + token)

	http.HandleFunc("/car", getdata.GetCarStatus)

	fmt.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}
}
