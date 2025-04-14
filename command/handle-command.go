package vehiclecommand

import (
	"context"
	"fmt"
	"time"

	"github.com/teslamotors/vehicle-command/pkg/account"
	"github.com/teslamotors/vehicle-command/pkg/protocol"
	"github.com/teslamotors/vehicle-command/pkg/vehicle"
)

func HandleCommand(req CommandRequest) error {
	if req.Vin == "" {
		return fmt.Errorf("no vehicle VIN found")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	privateKey, err := getPrivateKey()
	if err != nil {
		return fmt.Errorf("error getting private key: %s", err.Error())
	}

	userAgent := ""

	account, err := account.New(string(req.AuthToken), userAgent)
	if err != nil {
		return fmt.Errorf("error creating account: %s", err.Error())
	}

	car, err := account.GetVehicle(ctx, req.Vin, privateKey, nil)
	if err != nil {
		return fmt.Errorf("error getting vehicle: %s", err.Error())
	}

	err = car.Connect(ctx)
	if err != nil {
		return fmt.Errorf("error connecting to car: %s", err.Error())
	}

	err = car.StartSession(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting session: %s", err.Error())
	}

	err = handleIssueCommand(ctx, *car, req.Command)
	if err != nil {
		return fmt.Errorf("error issuing command: %s", err.Error())
	}

	return nil
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
	case "flash":
		err = car.FlashLights(ctx)
	default:
		err = fmt.Errorf("%s: is not a valid command", command)
	}
	if err != nil {
		return err
	}

	return nil
}

func getPrivateKey() (protocol.ECDHPrivateKey, error) {
	privateKey, err := protocol.LoadPrivateKey("./.temp/private-key.pem")
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}
