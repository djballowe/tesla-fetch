package main

import (
	"errors"
	"fmt"
	"sync"
	"tesla-app/client/api"
	"tesla-app/client/draw-status"
	"time"
)

type result struct {
	vehicleData api.VehicleData
	err         error
}

func main() {
	var group sync.WaitGroup
	group.Add(2)

	done := make(chan struct{})
	dataChan := make(chan result)

	go loading(&group, done)
	go getVehicleData(&group, done, dataChan)
	res := <-dataChan
	group.Wait()

	if res.err != nil {
		fmt.Println(res.err.Error())
		return
	}

	drawlogo.DrawStatus(res.vehicleData)

	return
}

func loading(group *sync.WaitGroup, done chan struct{}) {
	defer group.Done()
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-done:
			fmt.Printf("\r%s", "                         \n")
			return

		default:
			fmt.Printf("\r%s Fetching vehicle data", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}

func getVehicleData(group *sync.WaitGroup, done chan struct{}, dataChan chan result) {
	defer group.Done()
	defer close(done)
	carDataResponse, error := api.CallGetVehicleData()
	if error != nil {
		dataChan <- result{err: error}
		return
	}

	if carDataResponse.StatusCode == 401 {
		authResponse, error := api.CallAuth()
		if error != nil || authResponse.StatusCode != 200 {
			dataChan <- result{err: error}
			return
		}

		carDataResponse, error = api.CallGetVehicleData()
		if error != nil {
			dataChan <- result{err: error}
			return
		}
	}

	if carDataResponse.StatusCode != 200 {
		if carDataResponse.StatusCode == 408 {
			fmt.Println("Waking car")
			commandResp, error := api.CallIssueCommand("wake")
			if error != nil {
				dataChan <- result{err: error}
				return
			}

			if commandResp.StatusCode != 200 {
				error = errors.New(fmt.Sprintln("Could not issue command"))
				dataChan <- result{err: error}
				return
			}
		} else {
			error = errors.New(fmt.Sprintf("Error gathing vehicle data: Status Code %d", carDataResponse.StatusCode))
			dataChan <- result{err: error}
			return
		}
	}

	dataChan <- result{vehicleData: carDataResponse.Body, err: nil}
}
