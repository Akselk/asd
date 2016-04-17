package fsmModule

func stateMachine(state *int, event int, internalChan internalChanStruct) {

	switch *state {

	case IDLE:
		stateIdle(event, internalChan, state)

	case MOVING:
		stateMoving(event, internalChan, state)

	case DOOR_OPEN:
		stateDoorOpen(event, internalChan, state)

	}
}

func stateIdle(event int, internalChan internalChanStruct, state *int) {

	switch event {

	case START_MOVING_DOWN:
		*state = MOVING
		internalChan.goDown <- true

	case START_MOVING_UP:
		*state = MOVING
		internalChan.goUp <- true

	case REACHED_FLOOR:
		*state = DOOR_OPEN
		internalChan.openDoor <- true
		ExternalChan.TaskCompleted <- true
	}
}

func stateMoving(event int, internalChan internalChanStruct, state *int) {
	switch event {

	case REACHED_FLOOR:
		*state = DOOR_OPEN
		internalChan.stop <- true
		ExternalChan.TaskCompleted <- true
	}
}

func stateDoorOpen(event int, internalChan internalChanStruct, state *int) {
	switch event {
	case CLOSE_DOOR:
		*state = IDLE
		ExternalChan.StatusIdle <- true

	}

}
