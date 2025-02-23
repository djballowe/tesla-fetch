package vehiclecommand

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/vehicle"
)

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
		return handleReturn(fmt.Sprintf("Error getting vehicle: %s", err.Error()), false)
	}

	err = car.Connect(ctx)
	if err != nil {
		return handleReturn(fmt.Sprintf("Error connecting to car: %s", err.Error()), false)
	}

	err = car.StartSession(ctx, nil)
	if err != nil {
		return handleReturn(fmt.Sprintf("Error starting session: %s", err.Error()), false)
	}

	err = handleIssueCommand(ctx, *car, req.Command)
	if err != nil {
		return handleReturn(fmt.Sprintf("Error issuing command: %s", err.Error()), false)
	}

	success := fmt.Sprintf("Vehicle VIN: %s, command %s issued successfully", car.VIN(), req.Command)

	return handleReturn(success, true)

}

func handleIssueCommand(ctx context.Context, car vehicle.Vehicle, command string) error {
	var err error

	switch command {
	case "lock":
		err = car.Lock(ctx)
	case "unlock":
		err = car.Unlock(ctx)
	case "climate":
		err = car.ClimateOn(ctx)
	default:
		err = errors.New(fmt.Sprintf("%s: is not a valid command\n", command))
	}
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
