package data_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"tfetch/auth"
	"tfetch/data"
	"tfetch/ui"
	"tfetch/vehicle-state"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockVehicleService struct {
	VehicleStateFunc  func(token auth.Token) (*vehicle.VehicleStateResponse, error)
	PollWakeFunc      func(token auth.Token, status chan ui.ProgressUpdate) error
	VehicleStateCalls int
	PollWakeCalls     int
}

func (m *MockVehicleService) VehicleState(token auth.Token) (*vehicle.VehicleStateResponse, error) {
	if m.VehicleStateFunc != nil {
		m.VehicleStateCalls++
		return m.VehicleStateFunc(token)
	}

	panic("VehicleState mock not defined")
}

func (m *MockVehicleService) PollWake(token auth.Token, status chan ui.ProgressUpdate) error {
	if m.PollWakeFunc != nil {
		m.PollWakeCalls++
		return m.PollWakeFunc(token, status)
	}

	panic("PollWake mock not defined")
}

func (m *MockVehicleService) ResetMocks() {
	m.VehicleStateFunc = nil
	m.PollWakeFunc = nil
	m.VehicleStateCalls = 0
	m.PollWakeCalls = 0
}

func TestGetVehicleData(t *testing.T) {
	mockVehicleDataService := &MockVehicleService{}
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
		mockVehicleDataService.ResetMocks()

		mockVehicleDataService.VehicleStateFunc = func(token auth.Token) (*vehicle.VehicleStateResponse, error) {
			assert.Equal(t, mockToken, token)

			vehicleState := &vehicle.VehicleStateResponse{
				State: "online",
				Vin:   "",
			}
			return vehicleState, nil
		}
		mockVehicleDataService.PollWakeFunc = func(token auth.Token, status chan ui.ProgressUpdate) error {
			return nil
		}

		mockServerFunc := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/vehicles/%s/vehicle_data", testCarId), r.URL.Path)
			authHeader := r.Header.Get("Authorization")
			assert.Equal(t, "Bearer "+mockToken.AccessToken, authHeader)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := map[string]any{
				"response": map[string]any{
					"state": "online",
					"charge_state": map[string]any{
						"battery_level":          75,
						"battery_range":          133.6,
						"charge_rate":            10.0,
						"charging_state":         "Connected",
						"charge_limit_soc":       80,
						"minutes_to_full_charge": 10,
						"time_to_full_charge":    0.0,
					},
					"climate_state": map[string]any{
						"inside_temp":            21,
						"passenger_temp_setting": 22.0,
						"driver_temp_setting":    22.0,
						"is_climate_on":          false,
						"is_preconditioning":     false,
						"outside_temp":           18,
					},
					"vehicle_state": map[string]any{
						"locked":       true,
						"odometer":     12345,
						"vehicle_name": "Test Tesla Name",
					},
					"vehicle_config": map[string]any{
						"exterior_color":   "Black",
						"car_type":         "modely",
						"car_special_type": "base",
					},
				},
			}

			err := json.NewEncoder(w).Encode(response)
			require.NoError(t, err)
		}
		mockServer := httptest.NewServer(http.HandlerFunc(mockServerFunc))
		defer mockServer.Close()

		t.Setenv("TESLA_BASE_URL", mockServer.URL)

		var serviceMethods vehicle.VehicleMethods = mockVehicleDataService

		vehicleData, _ := data.GetVehicleData(mockStatus, mockToken, serviceMethods)

		// Should not call PollWake
		assert.Equal(t, 0, mockVehicleDataService.PollWakeCalls)
		assert.Equal(t, 1, mockVehicleDataService.VehicleStateCalls)

		// Vehicle data assertions
		assert.Equal(t, 75, vehicleData.BatteryLevel)
		assert.Equal(t, 133.6, vehicleData.BatteryRange)
		assert.Equal(t, 10.0, vehicleData.ChargeRate)
		assert.Equal(t, "Connected", vehicleData.ChargingState)
		assert.Equal(t, 80, vehicleData.ChargeLimitSoc)
		assert.Equal(t, 10, vehicleData.MinutesToFullCharge)
		assert.Equal(t, 0.0, vehicleData.TimeToFullCharge)
		assert.Equal(t, 21, vehicleData.InsideTemp)
		assert.Equal(t, 22.0, vehicleData.PassengerTempSetting)
		assert.Equal(t, 22.0, vehicleData.DriverTempSetting)
		assert.Equal(t, false, vehicleData.IsClimateOn)
		assert.Equal(t, false, vehicleData.IsPreconditioning)
		assert.Equal(t, 18, vehicleData.OutsideTemp)
		assert.Equal(t, true, vehicleData.Locked)
		assert.Equal(t, 12345, vehicleData.Odometer)
		assert.Equal(t, "Test Tesla Name", vehicleData.VehicleName)
		assert.Equal(t, "Black", vehicleData.Color)
		assert.Equal(t, "modely", vehicleData.CarType)
		assert.Equal(t, "base", vehicleData.CarSpecialType)
		assert.Equal(t, "online", vehicleData.State)

	})

	t.Run("Success - vehicle status offline", func(t *testing.T) {
		mockVehicleDataService.ResetMocks()

		mockVehicleDataService.VehicleStateFunc = func(token auth.Token) (*vehicle.VehicleStateResponse, error) {
			assert.Equal(t, mockToken, token)

			vehicleState := &vehicle.VehicleStateResponse{
				State: "offline",
				Vin:   "",
			}
			return vehicleState, nil
		}
		mockVehicleDataService.PollWakeFunc = func(token auth.Token, status chan ui.ProgressUpdate) error {
			return nil
		}

		mockServerFunc := func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "GET", r.Method)
			assert.Equal(t, fmt.Sprintf("/vehicles/%s/vehicle_data", testCarId), r.URL.Path)
			authHeader := r.Header.Get("Authorization")
			assert.Equal(t, "Bearer "+mockToken.AccessToken, authHeader)

			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			response := map[string]any{
				"response": map[string]any{
					"state": "online",
					"charge_state": map[string]any{
						"battery_level":          75,
						"battery_range":          133.6,
						"charge_rate":            10.0,
						"charging_state":         "Connected",
						"charge_limit_soc":       80,
						"minutes_to_full_charge": 10,
						"time_to_full_charge":    0.0,
					},
					"climate_state": map[string]any{
						"inside_temp":            21,
						"passenger_temp_setting": 22.0,
						"driver_temp_setting":    22.0,
						"is_climate_on":          false,
						"is_preconditioning":     false,
						"outside_temp":           18,
					},
					"vehicle_state": map[string]any{
						"locked":       true,
						"odometer":     12345,
						"vehicle_name": "Test Tesla Name",
					},
					"vehicle_config": map[string]any{
						"exterior_color":   "Black",
						"car_type":         "modely",
						"car_special_type": "base",
					},
				},
			}

			err := json.NewEncoder(w).Encode(response)
			require.NoError(t, err)
		}
		mockServer := httptest.NewServer(http.HandlerFunc(mockServerFunc))
		defer mockServer.Close()

		t.Setenv("TESLA_BASE_URL", mockServer.URL)

		var serviceMethods vehicle.VehicleMethods = mockVehicleDataService

		vehicleData, _ := data.GetVehicleData(mockStatus, mockToken, serviceMethods)

		// Should call PollWake
		assert.Equal(t, 1, mockVehicleDataService.PollWakeCalls)
		assert.Equal(t, 1, mockVehicleDataService.VehicleStateCalls)

		// Vehicle data assertions
		assert.Equal(t, 75, vehicleData.BatteryLevel)
		assert.Equal(t, 133.6, vehicleData.BatteryRange)
		assert.Equal(t, 10.0, vehicleData.ChargeRate)
		assert.Equal(t, "Connected", vehicleData.ChargingState)
		assert.Equal(t, 80, vehicleData.ChargeLimitSoc)
		assert.Equal(t, 10, vehicleData.MinutesToFullCharge)
		assert.Equal(t, 0.0, vehicleData.TimeToFullCharge)
		assert.Equal(t, 21, vehicleData.InsideTemp)
		assert.Equal(t, 22.0, vehicleData.PassengerTempSetting)
		assert.Equal(t, 22.0, vehicleData.DriverTempSetting)
		assert.Equal(t, false, vehicleData.IsClimateOn)
		assert.Equal(t, false, vehicleData.IsPreconditioning)
		assert.Equal(t, 18, vehicleData.OutsideTemp)
		assert.Equal(t, true, vehicleData.Locked)
		assert.Equal(t, 12345, vehicleData.Odometer)
		assert.Equal(t, "Test Tesla Name", vehicleData.VehicleName)
		assert.Equal(t, "Black", vehicleData.Color)
		assert.Equal(t, "modely", vehicleData.CarType)
		assert.Equal(t, "base", vehicleData.CarSpecialType)
		assert.Equal(t, "online", vehicleData.State)
	})
}
