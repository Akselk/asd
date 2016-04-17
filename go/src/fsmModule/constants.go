package fsmModule

import (
	"time"
)

const (
	DOOR_TIMER       = 3500 * time.Millisecond
	BREAK_DOWN_TIMER = 3000 * time.Millisecond
)

//states
const (
	IDLE      = 0
	MOVING    = 1
	DOOR_OPEN = 2
)

const (
	CLOSE_DOOR        = 4
	START_MOVING_UP   = 12
	START_MOVING_DOWN = 11
	REACHED_FLOOR     = 6
)
