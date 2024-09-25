package main

import (
	"fmt"
	//	"tesla-app/client/api"
	"sync"
	"tesla-app/client/draw-status"
	"time"
)

func main() {
	//	carDataResponse, error := api.CallGetVehicleData()
	//	if error != nil {
	//		fmt.Println(error.Error())
	//		return
	//	}
	//
	//	fmt.Println(carDataResponse.StatusCode)
	//
	//	if carDataResponse.StatusCode == 401 {
	//		authResponse, error := api.CallAuth()
	//		if error != nil || authResponse.StatusCode != 200 {
	//			fmt.Println(error.Error())
	//			return
	//		}
	//
	//		carDataResponse, error = api.CallGetVehicleData()
	//		if error != nil {
	//			fmt.Println(error.Error())
	//			return
	//		}
	//	}
	//
	//	if carDataResponse.StatusCode != 200 {
	//		if carDataResponse.StatusCode == 408 {
	//			fmt.Printf("Error gathering vehicle data: Status Code %d vehicle is asleep\n", carDataResponse.StatusCode)
	//			return
	//		}
	//		fmt.Printf("Error gathering vehicle data: Status Code %d\n", carDataResponse.StatusCode)
	//		return
	//	}

	var group sync.WaitGroup
	group.Add(2)

	flag := make(chan struct{})

	go loadingSpinner(&group, flag)
	go getVehicleData(&group, flag)

	group.Wait()

	drawlogo.DrawStatus()

	return
}

func loadingSpinner(group *sync.WaitGroup, flag chan struct{}) {
	defer group.Done()
	loadSpinner := [10]string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	idx := 0

	for {
		select {
		case <-flag:
			fmt.Printf("\r%s", "                         ")
			return

		default:
			fmt.Printf("\r%s Fetching vehicle data", loadSpinner[idx%10])
			time.Sleep(50 * time.Millisecond)
			idx++
		}
	}
}

func getVehicleData(group *sync.WaitGroup, flag chan struct{}) {
	defer group.Done()

	close(flag)
}
