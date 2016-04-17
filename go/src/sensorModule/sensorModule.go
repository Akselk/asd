package sensorModule

import (
	"colorPrint"
	"driverModule"
	"time"
)

func FloorSensor(floorChan chan int) {

	previousPosition := 10
	var currentPosition int

	for {

		time.Sleep(30 * time.Millisecond)

		currentPosition = driverModule.GetFloor()

		if (currentPosition != previousPosition) && (currentPosition != -1) {
			driverModule.SetFloorLight(currentPosition)
			colorPrint.DataWithColor(colorPrint.COLOR_BLUE, "Floor : ", currentPosition)

			select {

			case <-floorChan:
				floorChan <- currentPosition
			default:
				floorChan <- currentPosition
			}
			previousPosition = currentPosition
		}
	}

}
func ExternalButtonSensor(outsideButtonChan chan int) {

	for {
		time.Sleep(30 * time.Millisecond)
		for floorcounter := 0; floorcounter < driverModule.N_FLOORS; floorcounter++ {
			if driverModule.GetButton(driverModule.BUTTON_OUTSIDE_UP, floorcounter) {
				select {
				case <-outsideButtonChan:
					outsideButtonChan <- floorcounter
				default:
					outsideButtonChan <- floorcounter
				}
				time.Sleep(250 * time.Millisecond)
			}
			if driverModule.GetButton(driverModule.BUTTON_OUTSIDE_DOWN, floorcounter) {
				select {
				case <-outsideButtonChan:
					outsideButtonChan <- floorcounter + driverModule.N_FLOORS
				default:
					outsideButtonChan <- floorcounter + driverModule.N_FLOORS
				}
				time.Sleep(250 * time.Millisecond)
			}
		}
	}
}

func InternalButtonSensor(insideButtonChan chan int) {
	for {
		time.Sleep(30 * time.Millisecond)
		for floorcounter := 0; floorcounter < driverModule.N_FLOORS; floorcounter++ {
			if driverModule.GetButton(driverModule.BUTTON_INSIDE, floorcounter) {
				select {
				case <-insideButtonChan:
					insideButtonChan <- floorcounter
				default:
					insideButtonChan <- floorcounter
				}
				time.Sleep(250 * time.Millisecond)

			}

		}
	}
}
