package common

import (
	"context"
	"fmt"
	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/vehicle"
	"time"
)

type CommandRequest struct {
	AuthToken string `json:"authToken"`
	Vin       string `json:"vin"`
	Command   string `json:"command"`
}

type CommandResponse struct {
	Success bool   `json:"status"`
	Message string `json:"message"`
}

func HandleCommand(req CommandRequest) CommandResponse {
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

	err = handleIssueCommand(ctx, *car, req.Command)
	if err != nil {
		return handleReturn(err.Error(), false)
	}

	success := fmt.Sprintf("Vehicle VIN: %s, command %s issued successfully", car.VIN(), req.Command)

	return handleReturn(success, true)

}

func handleIssueCommand(ctx context.Context, car vehicle.Vehicle, command string) error {
	fmt.Println(command)
	//	err := car.HonkHorn(ctx)
	err := car.Wakeup(ctx)
	if err != nil {
		return err
	}
	return nil
}

func getPrivateKey() (protocol.ECDHPrivateKey, error) {
	privateKey, err := protocol.LoadPrivateKey("../.temp/private-key.pem")
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func handleReturn(err string, status bool) CommandResponse {
	return CommandResponse{
		Success: status,
		Message: err,
	}
}
