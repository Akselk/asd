//The systems har 3 queues.
// Internal for command orders.
//External for Call orders to be done by this elevator
// Global for all call orders from all elevators

// The queue is represented as a vector. The vector has one index for each command.
// For internal Queue that is the number of floors.
//For external and global queue it is twice that.

// Orders are added to the queue with a number valued one higher than the previous order.
// Thus the index in the queue represents the task to be done, and the value at this index
// decides how much this task is prioritized.

package queueModule

import (
	"backupModule"
	"driverModule"
	"fsmModule"
	"helpModule"
)

func RemoveFromGlobalQueue(Order int, GlobalQueue []int) {
	GlobalQueue[Order] = 0
}

func InsertToGlobalQueue(Order int, GlobalQueue []int) {
	if GlobalQueue[Order] == 0 {
		GlobalQueue[Order] = helpModule.FindBiggestNumber(GlobalQueue[:]) + 1

	}
}

func InsertToExternalQueue(buttonPressed int, External []int) {
	if External[buttonPressed] == 0 {
		External[buttonPressed] = helpModule.FindBiggestNumber(External[:]) + 1
	}
}

func InsertToInternalQueue(buttonPressed int, InternalQueue []int) {
	InternalQueue[buttonPressed] = helpModule.FindBiggestNumber(InternalQueue[:]) + 1
}

func PendingOrders(Queue []int) bool {
	if helpModule.FindBiggestNumber(Queue[:]) != 0 {
		return true
	}
	return false
}

func IsFirstTask(ExternalQueue []int, InternalQueue []int) bool {
	if helpModule.IsEmpty(ExternalQueue[:], InternalQueue[:]) {
		return true
	}
	return false
}

func IsNewTask(task int, internalQueue []int) bool {
	if internalQueue[task] == 0 {
		return true
	}
	return false
}

func SetInitialQueuePriority(queueEntry int, QueuePriority *int, floor int, QueueType string) {
	if QueueType == "Extern" {
		*QueuePriority = PRIORITY_EXTERNAL
	} else {
		if queueEntry >= floor {
			*QueuePriority = PRIORITY_INTERNAL_UP
		} else {
			*QueuePriority = PRIORITY_INTERNAL_DOWN
		}
	}
}

func OrderExecuted(ElevatorsInfo ElevatorsInfoStruct, Queue QueueStruct) (executedOrder int) {
	motorDir := driverModule.ReadAnalog(driverModule.MOTORDIR)

	DestinationFloor := helpModule.IndexToFloor(ElevatorsInfo.Destination[0])
	executedOrder = -1

	// Stopping for the active order
	if DestinationFloor == ElevatorsInfo.Floor[0] {
		driverModule.StopEngine()
		executedOrder = ElevatorsInfo.Destination[0]

		//Stopping for tasks on the way to the active order
	} else if Queue.External[ElevatorsInfo.Floor[0]+driverModule.N_FLOORS*motorDir] != 0 || Queue.Internal[ElevatorsInfo.Floor[0]] != 0 {
		driverModule.StopEngine()
		executedOrder = ElevatorsInfo.Floor[0] + driverModule.N_FLOORS*motorDir

	}

	return executedOrder
}

func MergeGlobalQueue(QuGlobal []int, oldGlobal []int) {
	for index := 0; index < 2*driverModule.N_FLOORS; index++ {
		if oldGlobal[index] != 0 {
			InsertToGlobalQueue(index, QuGlobal[:])
		}
	}
}

func StartNextOrder(Queue QueueStruct, ElevatorsInfo *ElevatorsInfoStruct, eventChan fsmModule.ExternalChanStruct, ExternalChan ExternalChanStruct) {

	switch Queue.Priority {

	case PRIORITY_EXTERNAL:
		ElevatorsInfo.Destination[0] = helpModule.FindIndexSmallestNumber(Queue.External[:])

	case PRIORITY_INTERNAL_UP:
		ElevatorsInfo.Destination[0] = (helpModule.HighestIndex(Queue.Internal[:]))

	case PRIORITY_INTERNAL_DOWN:
		ElevatorsInfo.Destination[0] = (helpModule.LowestIndex(Queue.Internal[:]))

	}
	DestinationFloor := helpModule.IndexToFloor(ElevatorsInfo.Destination[0])

	if ElevatorsInfo.Floor[0] < DestinationFloor {
		eventChan.StartMovingUp <- true
	} else if ElevatorsInfo.Floor[0] > DestinationFloor {
		eventChan.StartMovingDown <- true
	} else {
		ExternalChan.CheckFloor <- true
	}
}

func DeleteFromInternalQueue(internalQueue []int, order int) {
	floor := helpModule.IndexToFloor(order)
	internalQueue[floor] = 0
}

func DeleteFromExternalQueue(ExternalQueue []int, order int) {
	ExternalQueue[order] = 0
}
func InternalBackupQueueLoaded(internalQueue []int, floor int, QueuePriority *int) bool {
	backupQueue := backupModule.ReadInternalQuBackupFile()
	if PendingOrders(backupQueue[:]) {
		SetInitialQueuePriority(helpModule.FindIndexSmallestNumber(backupQueue[:]), QueuePriority, floor, "intern")
		for index := 0; index < len(backupQueue); index++ {
			internalQueue[index] = backupQueue[index]
			if internalQueue[index] != 0 {
				driverModule.ElevSetButtonLamp(driverModule.BUTTON_INSIDE, index, 1)
			}
		}
		return true
	}
	return false
}

func UpdateQueuePriority(Queue *QueueStruct, floor int, completedorder int) {

	switch Queue.Priority {

	case PRIORITY_EXTERNAL:
		if helpModule.FindBiggestNumber(Queue.Internal[:]) == 0 {
			Queue.Priority = PRIORITY_EXTERNAL
		} else if completedorder == helpModule.FindIndexSmallestNumber(Queue.External[:]) {
			if completedorder < driverModule.N_FLOORS {
				Queue.Priority = PRIORITY_INTERNAL_UP
			} else {
				Queue.Priority = PRIORITY_INTERNAL_DOWN
			}
		}

	case PRIORITY_INTERNAL_UP:
		completedorder = floor
		if helpModule.FindBiggestNumber(Queue.Internal[:]) == 0 {
			Queue.Priority = PRIORITY_EXTERNAL
		} else if completedorder == helpModule.HighestIndex(Queue.Internal[:]) {
			Queue.Priority = PRIORITY_INTERNAL_DOWN
		}

	case PRIORITY_INTERNAL_DOWN:
		completedorder = floor
		if helpModule.FindBiggestNumber(Queue.Internal[:]) == 0 {
			Queue.Priority = PRIORITY_EXTERNAL
		} else if completedorder == helpModule.LowestIndex(Queue.Internal[:]) {
			Queue.Priority = PRIORITY_INTERNAL_UP
		}
	}

}
