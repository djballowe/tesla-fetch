package data_test

import (
	"tesla-app/auth"
	"tesla-app/data"
	"tesla-app/ui"
	"tesla-app/vehicle-state"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockVehicleService struct {
	VehicleStateFunc func(token auth.Token) (*vehicle.VehicleStateResponse, error)
	PollWakeFunc     func(token auth.Token, status chan ui.ProgressUpdate) error
}

func (m *MockVehicleService) VehicleState(token auth.Token) (*vehicle.VehicleStateResponse, error) {
	if m.VehicleStateFunc != nil {
		return m.VehicleStateFunc(token)
	}

	panic("VehicleState mock not defined")
}

func (m *MockVehicleService) PollWake(token auth.Token, status chan ui.ProgressUpdate) error {
	if m.PollWakeFunc != nil {
		return m.PollWakeFunc(token, status)
	}

	panic("PollWake mock not defined")
}

func TestGetVehicleData(t *testing.T) {
	mockStatus := make(chan ui.ProgressUpdate, 1)
	mockToken := auth.Token{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		IdToken:      "test-id-token",
		State:        "test-state-token",
		TokenType:    "test-token-type",
		ExpiresIn:    1234,
		CreateAt:     time.Now(),
	}

	testCarId := "123456789"
	t.Setenv("MY_CAR_ID", testCarId)

	t.Run("Success - vehicle status online", func(t *testing.T) {
		mockVehicleDataService := &MockVehicleService{}
		mockVehicleDataService.VehicleStateFunc = func(token auth.Token) (*vehicle.VehicleStateResponse, error) {
			vehicleState := vehicle.VehicleStateResponse{
				State: "online",
				Vin:   "12345",
			}

			return &vehicleState, nil
		}
		mockVehicleDataService.PollWakeFunc = func(token auth.Token, status chan ui.ProgressUpdate) error {
			return nil
		}
		var serviceMethods vehicle.VehicleMethods = mockVehicleDataService

		vehicleData, _ := data.GetVehicleData(mockStatus, mockToken, serviceMethods)

		assert.Equal(t, "online", vehicleData.State)
	})
}
