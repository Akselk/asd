package sensorModule

type ExternalChanStruct struct {
	StopChan           chan int
	FloorChan          chan int
	ExternalButtonChan chan int
	InternalButtonChan chan int
}

func (sensorChan *ExternalChanStruct) Init() {

	sensorChan.StopChan = make(chan int, 1)
	sensorChan.FloorChan = make(chan int, 1)
	sensorChan.ExternalButtonChan = make(chan int, 1)
	sensorChan.InternalButtonChan = make(chan int, 1)
}
