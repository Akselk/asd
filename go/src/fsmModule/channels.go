package fsmModule

type internalChanStruct struct {
	openDoor            chan bool
	closeDoor           chan bool
	goUp                chan bool
	goDown              chan bool
	stop                chan bool
	startBreakdownTimer chan bool
	StopBreakdownTimer  chan bool
}

func (internalChan *internalChanStruct) initInternalFSMChan() {

	internalChan.openDoor = make(chan bool, 10)
	internalChan.closeDoor = make(chan bool, 10)
	internalChan.goUp = make(chan bool, 10)
	internalChan.goDown = make(chan bool, 10)
	internalChan.stop = make(chan bool, 10)
	internalChan.startBreakdownTimer = make(chan bool, 10)
	internalChan.StopBreakdownTimer = make(chan bool, 10)
}

type ExternalChanStruct struct {
	TaskCompleted         chan bool
	StopHere              chan bool
	StartMovingDown       chan bool
	StartMovingUp         chan bool
	Stop                  chan bool
	EngineError           chan bool
	GoDown                chan bool
	GoUp                  chan bool
	StatusIdle            chan bool
	ElevatorBreakdown     chan bool
	RestartBreakdownTimer chan bool
}

func (ExternalChan *ExternalChanStruct) Init() {

	ExternalChan.TaskCompleted = make(chan bool, 5)
	ExternalChan.StopHere = make(chan bool, 1)
	ExternalChan.StartMovingDown = make(chan bool, 1)
	ExternalChan.StartMovingUp = make(chan bool, 1)
	ExternalChan.Stop = make(chan bool, 1)
	ExternalChan.EngineError = make(chan bool, 1)
	ExternalChan.GoDown = make(chan bool, 2)
	ExternalChan.GoUp = make(chan bool, 2)
	ExternalChan.StatusIdle = make(chan bool, 1)
	ExternalChan.ElevatorBreakdown = make(chan bool, 1)
	ExternalChan.RestartBreakdownTimer = make(chan bool, 10)
}
