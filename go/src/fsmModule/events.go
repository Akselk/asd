package fsmModule

import (
	"driverModule"
	"fmt"
	"os"
	"time"
)

var ExternalChan ExternalChanStruct

func CheckEvents(ExternalChanInput ExternalChanStruct) {

	ExternalChan = ExternalChanInput

	var internalChan internalChanStruct
	internalChan.initInternalFSMChan()

	var state = IDLE

	for {
		time.Sleep(30 * time.Millisecond)

		select {

		case <-ExternalChan.StopHere:
			eventReachedFloor(&state, internalChan)

		case <-internalChan.stop:
			eventStop(&state, internalChan)

		case <-internalChan.goUp:
			eventGoUp(&state, internalChan)

		case <-internalChan.goDown:
			eventGoDown(&state, internalChan)

		case <-internalChan.openDoor:
			eventOpenDoor(&state, internalChan)

		case <-internalChan.closeDoor:
			eventCloseDoor(&state, internalChan, ExternalChan)

		case <-ExternalChan.StartMovingDown:
			eventStartMovingDown(&state, internalChan)

		case <-ExternalChan.StartMovingUp:
			eventStartMovingUp(&state, internalChan)

		case <-ExternalChan.Stop:
			eventStop(&state, internalChan)

		case <-ExternalChan.EngineError:
			eventEngineError(ExternalChan)

		case <-internalChan.startBreakdownTimer:
			go breakdownTimer(internalChan)

		}
	}
}

func eventReachedFloor(state *int, internalChan internalChanStruct) {
	stateMachine(state, REACHED_FLOOR, internalChan)
}

func eventStop(state *int, internalChan internalChanStruct) {
	driverModule.StopEngine()
	internalChan.openDoor <- true
	internalChan.StopBreakdownTimer <- true
}

func eventGoUp(state *int, internalChan internalChanStruct) {
	driverModule.StartEngine(driverModule.UP)
	driverModule.SetSpeed(2800)
}

func eventGoDown(state *int, internalChan internalChanStruct) {
	driverModule.StartEngine(driverModule.DOWN)
	driverModule.SetSpeed(-2800)
}

func eventOpenDoor(state *int, internalChan internalChanStruct) {
	driverModule.StopEngine()
	driverModule.SetDoorOpenLampON()
	fmt.Println("Opening door")
	go doorTimer(internalChan)
}

func eventCloseDoor(state *int, internalChan internalChanStruct, ExternalChan ExternalChanStruct) {
	driverModule.SetDoorOpenLampOFF()
	stateMachine(state, CLOSE_DOOR, internalChan)
}

func eventStartMovingUp(state *int, internalChan internalChanStruct) {
	stateMachine(state, START_MOVING_UP, internalChan)
	internalChan.startBreakdownTimer <- true
}

func eventStartMovingDown(state *int, internalChan internalChanStruct) {
	stateMachine(state, START_MOVING_DOWN, internalChan)
	internalChan.startBreakdownTimer <- true
}

func doorTimer(internalChan internalChanStruct) {
	select {
	case <-time.After(DOOR_TIMER):
		internalChan.closeDoor <- true
		fmt.Println("Closing door")
	}

}
func breakdownTimer(internalChan internalChanStruct) {
	select {
	case <-internalChan.StopBreakdownTimer:
		return

	case <-ExternalChan.RestartBreakdownTimer:
		go breakdownTimer(internalChan)
		return

	case <-time.After(BREAK_DOWN_TIMER):
		ExternalChan.EngineError <- true
	}
}

func eventEngineError(ExternalChan ExternalChanStruct) {
	driverModule.StopEngine()
	fmt.Println("===========================    ELEVATOR BREAKDOWN!   ===========================")
	fmt.Println("Elevator is obstructed or stuck  ->  Exiting program, restart imminent")

	os.Exit(1)
}
