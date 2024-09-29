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
	status := make(chan string)
	vehicleDataChan := make(chan api.VehicleData)

	go loading(&group, done, status)
	go getVehicleData(&group, done, status, vehicleDataChan)

	vehicleData := <-vehicleDataChan

	close(status)
	group.Wait()

	drawlogo.DrawStatus(vehicleData)

	return
}

func loading(group *sync.WaitGroup, done chan struct{}, status chan string) {
	defer group.Done()
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0
	msg := <-status

	for {
		select {
		case <-done:
			fmt.Printf("\r%s", "                         ")
			return

		default:
			fmt.Printf("\r%s %s", loadSpinner[idx%10], msg)
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}

func getVehicleData(group *sync.WaitGroup, done chan struct{}, status chan string, vehicleDataChan chan api.VehicleData) error {
	defer group.Done()
	carDataResponse, error := api.CallGetVehicleData()
	status <- "Data"
	if error != nil {
		return error
	}

	if carDataResponse.StatusCode == 401 {
		status <- "Auth"
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
			status <- "Asleep"
			return errors.New(fmt.Sprintf("Error gathering vehicle data: Status Code %d vehicle is asleep", carDataResponse.StatusCode))
		}
		return errors.New(fmt.Sprintf("Error gathing vehicle data: Status Code %d", carDataResponse.StatusCode))
	}

	vehicleDataChan <- carDataResponse.Body
	close(done)
	return nil
}
