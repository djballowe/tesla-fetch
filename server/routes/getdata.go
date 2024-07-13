package getdata

import (
	"fmt"
	"net/http"
	getTeslaAuth "tesla-app/server/common"
)

func GetCarStatus(writer http.ResponseWriter, req *http.Request) {
	tokenStore, state := getTeslaAuth.GetTokenStore()
	fmt.Println(state)
	fmt.Println("tokens: ", tokenStore[state].AccessToken)
	fmt.Println("here in get car status")

	if req.URL.Path != "/car" {
		http.Error(writer, "404 not found", http.StatusNotFound)
	}

	if req.Method != "GET" {
		http.Error(writer, "This method is not supported", http.StatusNotFound)
	}

	fmt.Fprintf(writer, "This is the car route")
}
