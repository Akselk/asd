package operatorModule

import (
	"backupModule"
	"costModule"
	"driverModule"
	"fsmModule"
	"helpModule"
	"networkModule"
	"queueModule"
	"sensorModule"
)

func ControlMainLoop(sensorChannels sensorModule.ExternalChanStruct, queueChannels queueModule.ExternalChanStruct, networkChannels networkModule.NetChannels, eventChannels fsmModule.ExternalChanStruct) {
	var (
		Queue         queueModule.QueueStruct
		ElevatorsInfo queueModule.ElevatorsInfoStruct
		localChan     localChannels
	)

	localChan.init()
	ElevatorsInfo.ElevatorsInfoInit()
	driverModule.ElevatorInit()
	if queueModule.InternalBackupQueueLoaded(Queue.Internal[:], ElevatorsInfo.Floor[0], &Queue.Priority) {
		localChan.findNextTask <- true
	}

	for {
		select {

		case ip := <-networkChannels.NewConnection:
			ElevatorsInfo = appendNewElevator(ElevatorsInfo, ip)
			handleNewConnection(ElevatorsInfo, Queue, localChan.outbox)
			continue

		case Inbox := <-networkChannels.Inbox:
			handleRecievedMail(&ElevatorsInfo, &Queue, Inbox, localChan)
			continue

		case outbox := <-localChan.outbox:
			handleSendMail(outbox, ElevatorsInfo.IP[:], networkChannels)
			continue

		case floor := <-sensorChannels.FloorChan:
			handleNewFloor(floor, ElevatorsInfo, eventChannels.RestartBreakdownTimer, localChan.outbox, queueChannels.CheckFloor, Queue)

		case ExternalButton := <-sensorChannels.ExternalButtonChan:
			handleNewExternalOrder(ExternalButton, ElevatorsInfo, Queue, localChan.outbox)

		case <-queueChannels.CheckFloor:
			if shouldStopAtFloor(ElevatorsInfo, Queue) {
				eventChannels.StopHere <- true
			}

		case <-eventChannels.TaskCompleted:
			handleCompletedTask(ElevatorsInfo, &Queue, localChan.outbox)

		case internalButton := <-sensorChannels.InternalButtonChan:
			handleNewInternalOrder(internalButton, ElevatorsInfo.Floor[0], &Queue, localChan.findNextTask)

		case <-eventChannels.StatusIdle:
			setIdle(ElevatorsInfo.Destination[:])
			localChan.outbox <- CreateMail(UPDATE_DESTINATION, ElevatorsInfo, Queue, NOT_AN_ORDER)
			localChan.findNextTask <- true

		case <-localChan.findNextTask:
			if costModule.OptimalTaskFound(ElevatorsInfo, &Queue) {
				queueModule.StartNextOrder(Queue, &ElevatorsInfo, eventChannels, queueChannels)
				localChan.outbox <- CreateMail(UPDATE_DESTINATION, ElevatorsInfo, Queue, NOT_AN_ORDER)
			}

		case deadElevIP := <-networkChannels.GetDeadElevator:
			ElevatorsInfo = removeDeadElevatorInfo(deadElevIP, ElevatorsInfo)
			localChan.findNextTask <- true
			continue

		}
	}
}

func handleNewConnection(ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue queueModule.QueueStruct, outbox chan networkModule.Mail) {
	outbox <- CreateMail(UPDATE_DESTINATION, ElevatorsInfo, Queue, NOT_AN_ORDER)
	outbox <- CreateMail(UPDATE_FLOOR, ElevatorsInfo, Queue, NOT_AN_ORDER)
	if queueModule.PendingOrders(Queue.Global[:]) {
		outbox <- CreateMail(MERGE_GLOBAL_QUEUE, ElevatorsInfo, Queue, 0)
	}
}

func appendNewElevator(ElevatorsInfo queueModule.ElevatorsInfoStruct, IP string) (ElevatorsInfoOut queueModule.ElevatorsInfoStruct) {
	ElevatorsInfoOut.IP = append(ElevatorsInfo.IP, IP)
	ElevatorsInfoOut.Floor = append(ElevatorsInfo.Floor, -2)
	ElevatorsInfoOut.Destination = append(ElevatorsInfo.Destination, -2)
	return ElevatorsInfoOut
}

func handleRecievedMail(ElevatorsInfo *queueModule.ElevatorsInfoStruct, Queue *queueModule.QueueStruct, Inbox networkModule.Mail, localChan localChannels) {
	switch Inbox.Msg.Type {

	case NEW_GLOBAL_QUEUE:
		Queue.Global = Inbox.Msg.GlobalQueue
		driverModule.SetExternalOrderLamps(Queue.Global[:])
		localChan.findNextTask <- true

	case MERGE_GLOBAL_QUEUE:
		if helpModule.FindLowestIPindex(ElevatorsInfo.IP[:]) == 0 {
			queueModule.MergeGlobalQueue(Queue.Global[:], Inbox.Msg.GlobalQueue[:])
			localChan.outbox <- CreateMail(NEW_GLOBAL_QUEUE, *ElevatorsInfo, *Queue, NOT_AN_ORDER)
			localChan.findNextTask <- true
		}

	case UPDATE_DESTINATION:
		updateElevatorsInfoDestination(ElevatorsInfo, Inbox)
		localChan.findNextTask <- true

	case UPDATE_FLOOR:
		updateElevatorsInfoFloors(ElevatorsInfo, Inbox)

	case ADD_ORDER:
		queueModule.InsertToGlobalQueue(Inbox.Msg.Order, Queue.Global[:])
		localChan.outbox <- CreateMail(NEW_GLOBAL_QUEUE, *ElevatorsInfo, *Queue, NOT_AN_ORDER)
		localChan.findNextTask <- true
		driverModule.SetExternalOrderLamps(Queue.Global[:])

	case REMOVE_ORDER:
		queueModule.RemoveFromGlobalQueue(Inbox.Msg.Order, Queue.Global[:])
		localChan.outbox <- CreateMail(NEW_GLOBAL_QUEUE, *ElevatorsInfo, *Queue, NOT_AN_ORDER)
		driverModule.SetExternalOrderLamps(Queue.Global[:])

	}
}

func handleNewFloor(Floor int, ElevatorsInfo queueModule.ElevatorsInfoStruct, restartBreakdownTimer chan bool, outbox chan networkModule.Mail, checkFloor chan bool, Queue queueModule.QueueStruct) {
	ElevatorsInfo.Floor[0] = Floor
	restartBreakdownTimer <- true
	outbox <- CreateMail(UPDATE_FLOOR, ElevatorsInfo, Queue, NOT_AN_ORDER)
	checkFloor <- true
}

func handleNewExternalOrder(ExternalButton int, ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue queueModule.QueueStruct, outbox chan networkModule.Mail) {
	if Queue.Global[ExternalButton] == 0 {
		outbox <- CreateMail(ADD_ORDER, ElevatorsInfo, Queue, ExternalButton)
	}
}

func CreateMail(Type int, ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue queueModule.QueueStruct, order int) (mail networkModule.Mail) {
	var IP string
	if order == NOT_AN_ORDER {
		IP = ""
	} else {
		IP = "LowestIP"
	}
	mail.MakeMail(IP, ElevatorsInfo.IP[0], Type, ElevatorsInfo.Destination[0], ElevatorsInfo.Floor[0], order, Queue.Global)
	return
}

func handleSendMail(Outbox networkModule.Mail, ElevatorsIP []string, netChan networkModule.NetChannels) {
	if Outbox.TargetIP == "" {
		netChan.SendToAll <- Outbox
	} else {
		Outbox.TargetIP = ElevatorsIP[helpModule.FindLowestIPindex(ElevatorsIP[:])]
	}
	if Outbox.TargetIP == ElevatorsIP[0] {
		netChan.Inbox <- Outbox
	} else {
		netChan.SendToOne <- Outbox
	}
}

func handleNewInternalOrder(insideElevButton int, floor int, Queue *queueModule.QueueStruct, findNextTask chan bool) {
	if queueModule.IsNewTask(insideElevButton, Queue.Internal[:]) {
		if queueModule.IsFirstTask(Queue.External[:], Queue.Internal[:]) {
			queueModule.SetInitialQueuePriority(insideElevButton, &Queue.Priority, floor, "intern")
			findNextTask <- true
		}
		queueModule.InsertToInternalQueue(insideElevButton, Queue.Internal[:])
		driverModule.TurnOnInternalButtonLight(insideElevButton)
		backupModule.WriteInternalQueueBackupFile(Queue.Internal[:])
	}
}

func handleCompletedTask(ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue *queueModule.QueueStruct, outbox chan networkModule.Mail) {
	executedOrder := queueModule.OrderExecuted(ElevatorsInfo, *Queue)
	queueModule.DeleteFromInternalQueue(Queue.Internal[:], executedOrder)
	backupModule.WriteInternalQueueBackupFile(Queue.Internal[:])
	queueModule.UpdateQueuePriority(Queue, ElevatorsInfo.Floor[0], executedOrder)
	driverModule.TurnOffInternalButtonLight(executedOrder)
	queueModule.DeleteFromExternalQueue(Queue.External[:], executedOrder)
	outbox <- CreateMail(REMOVE_ORDER, ElevatorsInfo, *Queue, executedOrder)
}

func updateElevatorsInfoFloors(ElevatorsInfo *queueModule.ElevatorsInfoStruct, mail networkModule.Mail) {
	for index := 1; index < len(ElevatorsInfo.IP); index++ {
		if mail.Msg.SendersIP == ElevatorsInfo.IP[index] {
			ElevatorsInfo.Floor[index] = mail.Msg.Floor
		}
	}
}

func updateElevatorsInfoDestination(ElevatorsInfo *queueModule.ElevatorsInfoStruct, mail networkModule.Mail) {
	for index := 1; index < len(ElevatorsInfo.IP); index++ {
		if mail.Msg.SendersIP == ElevatorsInfo.IP[index] {
			ElevatorsInfo.Destination[index] = mail.Msg.Destination
		}
	}
}

func shouldStopAtFloor(ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue queueModule.QueueStruct) bool {
	if queueModule.OrderExecuted(ElevatorsInfo, Queue) == -1 {
		return false
	}
	return true
}

func setIdle(destination []int) {
	destination[0] = -1
}

func removeDeadElevatorInfo(deadIP string, ElevatorsInfo queueModule.ElevatorsInfoStruct) (ElevatorsInfoOut queueModule.ElevatorsInfoStruct) {
	for index := 0; index < len(ElevatorsInfo.IP); index++ {
		if ElevatorsInfo.IP[index] == deadIP {
			ElevatorsInfo.IP = append(ElevatorsInfo.IP[:index], ElevatorsInfo.IP[index+1:]...)
			ElevatorsInfo.Floor = append(ElevatorsInfo.Floor[:index], ElevatorsInfo.Floor[index+1:]...)
			ElevatorsInfo.Destination = append(ElevatorsInfo.Destination[:index], ElevatorsInfo.Destination[index+1:]...)
			break

		}
	}

	return ElevatorsInfo
}
