package routes

import (
	"fmt"
	"net/http"
	"os"
	"tesla-app/server/common"
)

func IssueCommand(writer http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		http.Error(writer, "This method is not supported", http.StatusNotFound)
	}

	req.ParseForm()
	command := req.Form.Get("command")

	tokenStore, state := common.GetTokenStore()
	vin := os.Getenv("VIN")

	authToken := tokenStore[state].AccessToken
	var commandReq = common.CommandRequest{
		AuthToken: authToken,
		Vin:       vin,
		Command:   command,
	}

	commandResp := common.HandleCommand(commandReq)

	fmt.Printf("Command issue status: Success? %t Message: %s", commandResp.Success, commandResp.Message)
}
