package getdata

import (
	"fmt"
	"net/http"
)

func GetCarStatus(writer http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/car" {
		http.Error(writer, "404 not found", http.StatusNotFound)
	}

	if req.Method != "GET" {
		http.Error(writer, "This method is not supported", http.StatusNotFound)
	}

	fmt.Fprintf(writer, "This is the car route")
}
