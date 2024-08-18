package common

import (
	"context"
	"fmt"
	"time"

	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
)

type Request struct {
	AuthToken string `json:"authToken"`
	Vin       string `json:"vin"`
	Command   string `json:"command"`
}

type Response struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}

func HandleCommand(req Request) {
	if req.Vin == "" {
		fmt.Println("No vehicle VIN found")
		//		return Response{
		//			Success: false,
		//			Message: "No vehicle VIN found",
		//		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	privateKey, err := getPrivateKey()
	if err != nil {
		fmt.Println(err.Error())
		//		return Response{
		//			Success: false,
		//			Message: "Could not load private key",
		//		}
	}

	userAgent := ""

	account, err := account.New(string(req.AuthToken), userAgent)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	car, err := account.GetVehicle(ctx, req.Vin, privateKey, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("VIN: %s\n", car.VIN())
	fmt.Println("Connecting to car...")

	err = car.Connect(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
		//		return Response{
		//			Success: false,
		//			Message:"Could not connect to car",
		//		}
	}

	err = car.StartSession(ctx, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = car.HonkHorn(ctx)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}

func getPrivateKey() (protocol.ECDHPrivateKey, error) {
	privateKey, err := protocol.LoadPrivateKey("../.temp/private-key.pem")
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
