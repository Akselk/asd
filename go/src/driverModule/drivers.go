package driverModule

import (
	"fmt"
	"helpModule"
	"time"
)

const N_BUTTONS = 3
const N_FLOORS = 4
const (
	UP    = 0x01
	DOWN  = 0x02
	STILL = 0x03
)

var lampChannelMatrix = [N_FLOORS][N_BUTTONS]int{
	{LIGHT_UP1, LIGHT_DOWN1, LIGHT_COMMAND1},
	{LIGHT_UP2, LIGHT_DOWN2, LIGHT_COMMAND2},
	{LIGHT_UP3, LIGHT_DOWN3, LIGHT_COMMAND3},
	{LIGHT_UP4, LIGHT_DOWN4, LIGHT_COMMAND4}}

var buttonChannelMatrix = [N_FLOORS][N_BUTTONS]int{
	{BUTTON_UP1, BUTTON_DOWN1, BUTTON_COMMAND1},
	{BUTTON_UP2, BUTTON_DOWN2, BUTTON_COMMAND2},
	{BUTTON_UP3, BUTTON_DOWN3, BUTTON_COMMAND3},
	{BUTTON_UP4, BUTTON_DOWN4, BUTTON_COMMAND4}}

type ButtonType int

const (
	BUTTON_OUTSIDE_UP   = 0
	BUTTON_OUTSIDE_DOWN = 1
	BUTTON_INSIDE       = 2
)

func ElevatorInit() {

	Init(ET_comedi)
	TurnOffLights()
	StopEngine()

	return
}

func DriveToClosestFloor() {

	if GetFloor() != -1 {
		return
	}

	once := true
	ElevatorInit()
	StartEngine(UP)
	SetSpeed(2800)

	for timer := 0; timer < 5000; timer += 50 {
		if GetFloor() != -1 {
			StopEngine()
			return
		} else if timer > 2500 && once {
			StopEngine()
			time.Sleep(100 * time.Millisecond)
			StartEngine(DOWN)
			SetSpeed(-2800)
			once = false
		}
		time.Sleep(50 * time.Millisecond)
	}

	StopEngine()
	fmt.Println("TIMED OUT : Did not find any initial floor placement -> Check elevator for obstruction")
}

func GetFloor() int {

	if ReadBit(SENSOR_FLOOR1) {
		return 0
	} else if ReadBit(SENSOR_FLOOR2) {
		return 1
	} else if ReadBit(SENSOR_FLOOR3) {
		return 2
	} else if ReadBit(SENSOR_FLOOR4) {
		return 3
	} else {
		return -1
	}
}

func GetButton(button ButtonType, floor int) bool {

	if ReadBit(buttonChannelMatrix[floor][button]) {
		return true
	} else {
		return false
	}
}

func SetSpeed(speed int) {

	if speed > 0 {
		ClearBit(MOTORDIR)
	} else if speed < 0 {
		SetBit(MOTORDIR)
		speed = -speed
	}

	WriteAnalog(MOTOR, speed)
}

func StartEngine(direction int) {

	if direction == UP {
		ClearBit(MOTORDIR)
	} else {
		SetBit(MOTORDIR)
	}

	WriteAnalog(MOTOR, MOTOR_SPEED)
}

func StopEngine() {
	WriteAnalog(MOTOR, 0)
}

func SetDoorOpenLampON() {
	SetBit(LIGHT_DOOR_OPEN)
}

func SetDoorOpenLampOFF() {
	ClearBit(LIGHT_DOOR_OPEN)
}

func GetStopSignal() bool {
	return ReadBit(STOP)
}

func SetFloorLight(floor int) {

	if (floor & 0x02) != 0 {
		SetBit(LIGHT_FLOOR_IND1)
	} else {
		ClearBit(LIGHT_FLOOR_IND1)
	}
	if (floor & 0x01) != 0 {
		SetBit(LIGHT_FLOOR_IND2)
	} else {
		ClearBit(LIGHT_FLOOR_IND2)
	}
}

func ElevSetButtonLamp(button ButtonType, floor int, value int) {
	if floor >= 0 && floor <= N_FLOORS {

		if value == 1 {
			SetBit(lampChannelMatrix[floor][button])
		} else {
			ClearBit(lampChannelMatrix[floor][button])
		}
	}
}

func TurnOffInternalButtonLight(order int) {
	floor := helpModule.IndexToFloor(order)
	if floor >= 0 && floor <= N_FLOORS {
		ClearBit(lampChannelMatrix[floor][BUTTON_INSIDE])
	}
}

func TurnOnInternalButtonLight(order int) {
	floor := helpModule.IndexToFloor(order)
	if floor >= 0 && floor <= N_FLOORS {
		SetBit(lampChannelMatrix[floor][BUTTON_INSIDE])
	}
}

func SetExternalOrderLamps(Queue []int) {
	for index := 0; index < len(Queue)/2; index++ {
		if Queue[index] != 0 {
			ElevSetButtonLamp(BUTTON_OUTSIDE_UP, index, 1)
		} else {
			ElevSetButtonLamp(BUTTON_OUTSIDE_UP, index, 0)
		}
		if Queue[index+N_FLOORS] != 0 {
			ElevSetButtonLamp(BUTTON_OUTSIDE_DOWN, index, 1)
		} else {
			ElevSetButtonLamp(BUTTON_OUTSIDE_DOWN, index, 0)
		}
	}
}

func TurnOffLights() {
	for floors := 0; floors < N_FLOORS; floors++ {
		ElevSetButtonLamp(0, floors, 0)
		ElevSetButtonLamp(1, floors, 0)
		ElevSetButtonLamp(2, floors, 0)

	}
}
