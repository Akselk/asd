package operatorModule

import "networkModule"

type localChannels struct {
	triggerOrder    chan int
	orderCompleted  chan bool
	findNextTask    chan bool
	executeNextTask chan bool
	outbox          chan networkModule.Mail
}

func (localChan *localChannels) init() {
	localChan.triggerOrder = make(chan int, 12)
	localChan.orderCompleted = make(chan bool, 12)
	localChan.executeNextTask = make(chan bool, 12)
	localChan.findNextTask = make(chan bool, 12)
	localChan.outbox = make(chan networkModule.Mail, 20)
}
