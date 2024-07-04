package main

import (
	"fmt"
	"log"
	"net/http"
	"tesla-app/server/routes"
)

func main() {
	http.HandleFunc("/car", getdata.GetCarStatus)

	fmt.Println("Starting server on port: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err == nil {
		log.Fatal(err)
	}

}
