package queueModule

import (
	"driverModule"
	"networkModule"
)

type QueueStruct struct {
	Global   [driverModule.N_FLOORS * 2]int
	External [driverModule.N_FLOORS * 2]int
	Internal [driverModule.N_FLOORS]int
	Priority int
}

type ElevatorsInfoStruct struct {
	IP          []string
	Destination []int
	Floor       []int
}

func (ElevatorsInfo *ElevatorsInfoStruct) ElevatorsInfoInit() {
	ElevatorsInfo.IP = []string{networkModule.GetLocalIP()}
	ElevatorsInfo.Destination = []int{-1}
	ElevatorsInfo.Floor = []int{driverModule.GetFloor()}
}

const (
	PRIORITY_EXTERNAL      = 1
	PRIORITY_INTERNAL_UP   = 2
	PRIORITY_INTERNAL_DOWN = 3
)
