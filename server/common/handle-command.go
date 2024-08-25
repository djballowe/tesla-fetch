package common

import (
	"context"
	"fmt"
	"time"

	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
//	"github.com/teslamotors/vehicle-command/pkg/vehicle"
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

func HandleCommand(req Request) Response {
	if req.Vin == "" {
		return handleReturn("No vehicle VIN found", false)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	privateKey, err := getPrivateKey()
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	userAgent := ""

	account, err := account.New(string(req.AuthToken), userAgent)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	car, err := account.GetVehicle(ctx, req.Vin, privateKey, nil)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	fmt.Printf("VIN: %s\n", car.VIN())
	fmt.Println("Connecting to car...")

	err = car.Connect(ctx)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	err = car.StartSession(ctx, nil)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	err = car.HonkHorn(ctx)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	//	err = handleIssueCommand(ctx, *car, req.Command)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	success := fmt.Sprintf("Vehicle VIN: %s, command %s issued successfully", car.VIN(), req.Command)

	return handleReturn(success, true)

}

func getPrivateKey() (protocol.ECDHPrivateKey, error) {
	privateKey, err := protocol.LoadPrivateKey("../.temp/private-key.pem")
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func handleReturn(err string, status bool) Response {
	return Response{
		Success: status,
		Message: err,
	}
}

//func handleIssueCommand(ctx context.Context, vehicle vehicle.Vehicle, action string) error {
//	return nil
//}
