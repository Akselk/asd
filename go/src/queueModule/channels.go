package queueModule

type ExternalChanStruct struct {
	CheckFloor chan bool
}

func (queueChan *ExternalChanStruct) Init() {
	queueChan.CheckFloor = make(chan bool, 12)
}
