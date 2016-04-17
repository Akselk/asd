package costModule

import (
	"helpModule"
	"queueModule"
)

func OptimalTaskFound(ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue *queueModule.QueueStruct) bool {

	oldExternalQueue := Queue.External
	CalculateMyExternalQueue(ElevatorsInfo, Queue)
	if oldExternalQueue != Queue.External {
		if queueModule.IsFirstTask(oldExternalQueue[:], Queue.Internal[:]) {
			queueModule.SetInitialQueuePriority(helpModule.FindIndexSmallestNumber(Queue.External[:]), &Queue.Priority, ElevatorsInfo.Floor[0], "Extern")
		}
	}

	if !helpModule.IsEmpty(Queue.External[:], Queue.Internal[:]) && ElevatorsInfo.Destination[0] == -1 {
		return true
	}
	return false
}

func CalculateMyExternalQueue(ElevatorsInfo queueModule.ElevatorsInfoStruct, Queue *queueModule.QueueStruct) {
	if queueModule.PendingOrders(Queue.Global[:]) {
		ActiveQueueEntry := helpModule.FindIndexSmallestNumber(Queue.Global[:])

		for ActiveQueueEntry != -1 {

			if !isElevatorsonTask(ElevatorsInfo.Destination[:], ActiveQueueEntry) {

				if isIdleElevators(ElevatorsInfo.Destination[:]) {
					IndexBestElev := caluclateRightElevator(ElevatorsInfo, ActiveQueueEntry)

					if IndexBestElev == 0 {
						if Queue.External[ActiveQueueEntry] == 0 {
							queueModule.InsertToExternalQueue(ActiveQueueEntry, Queue.External[:])
						}
					} else {
						Queue.External[ActiveQueueEntry] = 0
					}
				} else {
					IndexBestElev := calculateDriveByElev(ActiveQueueEntry, ElevatorsInfo)
					if IndexBestElev == 0 {
						queueModule.InsertToExternalQueue(ActiveQueueEntry, Queue.External[:])
					}
				}

			}
			updateActiveQueueEntry(Queue.Global[:], &ActiveQueueEntry)
		}

	}
}

func isIdleElevators(CurrDest []int) bool {
	if helpModule.FindIndexOfNumber(CurrDest[:], -1) != -1 {
		return true
	} else {
		return false
	}

}

func isElevatorsonTask(CurrDest []int, ActiveQueueEntry int) bool {
	if helpModule.FindIndexOfNumber(CurrDest[:], ActiveQueueEntry) != -1 {
		return true
	} else {
		return false
	}
}

func caluclateRightElevator(ElevatorsInfo queueModule.ElevatorsInfoStruct, ActiveQueueEntry int) (BestElev int) {
	lowestDistanceToMission := 100
	BestElev = -1
	DistancetoMission := 100
	lowestIP := "999.999.999.999"

	for elevIndex := 0; elevIndex < len(ElevatorsInfo.Destination); elevIndex++ {
		if ElevatorsInfo.Destination[elevIndex] == -1 || ElevatorsInfo.Destination[elevIndex] == ActiveQueueEntry {
			DistancetoMission = helpModule.Abs(ElevatorsInfo.Floor[elevIndex] - helpModule.IndexToFloor(ActiveQueueEntry))

			if DistancetoMission < lowestDistanceToMission {
				lowestDistanceToMission = DistancetoMission
				lowestIP = ElevatorsInfo.IP[elevIndex]
				BestElev = elevIndex
			} else if DistancetoMission == lowestDistanceToMission {
				if ElevatorsInfo.IP[elevIndex] < lowestIP {
					lowestIP = ElevatorsInfo.IP[elevIndex]
					BestElev = elevIndex
				}

			}

		}
	}
	return
}

func calculateDriveByElev(ActiveQueueEntry int, ElevatorsInfo queueModule.ElevatorsInfoStruct) (BestElev int) {
	lowestDistanceToMission := 100
	BestElev = -1
	DistancetoMission := 100
	lowestIP := "999.999.999.999"

	DriveByDestinationFloor := helpModule.IndexToFloor(ActiveQueueEntry)
	ElevatorDestinationFloor := -1

	for elevIndex := 0; elevIndex < len(ElevatorsInfo.Destination); elevIndex++ {
		ElevatorDestinationFloor = helpModule.IndexToFloor(ElevatorsInfo.Destination[elevIndex])
		if (DriveByDestinationFloor > ElevatorsInfo.Floor[elevIndex] && DriveByDestinationFloor < ElevatorDestinationFloor) || (DriveByDestinationFloor < ElevatorsInfo.Floor[elevIndex] && DriveByDestinationFloor > ElevatorDestinationFloor) {
			DistancetoMission = helpModule.Abs(ElevatorsInfo.Floor[elevIndex] - ActiveQueueEntry)

			if DistancetoMission < lowestDistanceToMission {
				lowestDistanceToMission = DistancetoMission
				lowestIP = ElevatorsInfo.IP[elevIndex]
				BestElev = elevIndex
			} else if DistancetoMission == lowestDistanceToMission {
				if ElevatorsInfo.IP[elevIndex] < lowestIP {
					lowestIP = ElevatorsInfo.IP[elevIndex]
					BestElev = elevIndex
				}

			}

		}
	}
	return
}

func updateActiveQueueEntry(GloQu []int, ActiveQueueEntry *int) {
	*ActiveQueueEntry = helpModule.FindIndexSmallestNumberLargerThan(GloQu[:], GloQu[*ActiveQueueEntry])
}
