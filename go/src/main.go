package main

import (
	"driverModule"
	"fsmModule"
	"networkModule"
	"operatorModule"
	"queueModule"
	"sensorModule"
)

func main() {

	var (
		networkChannels networkModule.NetChannels
		sensorChannels  sensorModule.ExternalChanStruct
		queueChannels   queueModule.ExternalChanStruct
		eventChannels   fsmModule.ExternalChanStruct
	)

	networkChannels.Init()
	sensorChannels.Init()
	queueChannels.Init()
	eventChannels.Init()

	driverModule.DriveToClosestFloor()

	go networkModule.NetworkSetup(networkChannels)
	go fsmModule.CheckEvents(eventChannels)

	go sensorModule.ExternalButtonSensor(sensorChannels.ExternalButtonChan)
	go sensorModule.InternalButtonSensor(sensorChannels.InternalButtonChan)
	go sensorModule.FloorSensor(sensorChannels.FloorChan)
	operatorModule.ControlMainLoop(sensorChannels, queueChannels, networkChannels, eventChannels)
}
