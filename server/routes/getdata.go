package getdata

import (
	"fmt"
	"net/http"
//	getTeslaAuth "tesla-app/server/common"
)

func GetCarStatus(writer http.ResponseWriter, req *http.Request) {
//	store, tokenStore, stateStore := getTeslaAuth.GetTokenStore()
//	store.Lock()
//	fmt.Println("store: ", store)
//	fmt.Println("store: ", *tokenStore.AccessToken)
//	store.Unlock()

	if req.URL.Path != "/car" {
		http.Error(writer, "404 not found", http.StatusNotFound)
	}

	if req.Method != "GET" {
		http.Error(writer, "This method is not supported", http.StatusNotFound)
	}

	fmt.Fprintf(writer, "This is the car route")
}
