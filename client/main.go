package main

import (
	"errors"
	"fmt"
	"sync"
	"tesla-app/client/api"
	"tesla-app/client/draw-status"
	"time"
)

func main() {

	var group sync.WaitGroup
	group.Add(2)

	done := make(chan struct{})
	vehicleDataChan := make(chan api.VehicleData)

	go loadingSpinner(&group, done)
	go getVehicleData(&group, done, vehicleDataChan)
	vehicleData := <-vehicleDataChan

	group.Wait()

	drawlogo.DrawStatus(vehicleData)

	return
}

func loadingSpinner(group *sync.WaitGroup, done chan struct{}) {
	defer group.Done()
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s", "                         ")
			return

		default:
			fmt.Printf("\r%s Fetching vehicle data", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}

func getVehicleData(group *sync.WaitGroup, done chan struct{}, vehicleDataChan chan api.VehicleData) error {
	defer group.Done()
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		return error
	}

	if carDataResponse.StatusCode == 401 {
		authResponse, error := api.CallAuth()
		if error != nil || authResponse.StatusCode != 200 {
			return error
		}

		carDataResponse, error = api.CallGetVehicleData()
		if error != nil {
			return error
		}
	}

	if carDataResponse.StatusCode != 200 {
		if carDataResponse.StatusCode == 408 {
			return errors.New(fmt.Sprintf("Error gathering vehicle data: Status Code %d vehicle is asleep", carDataResponse.StatusCode))
		}
		return errors.New(fmt.Sprintf("Error gathing vehicle data: Status Code %d", carDataResponse.StatusCode))
	}

	vehicleDataChan <- carDataResponse.Body
	close(done)
	return nil
}
