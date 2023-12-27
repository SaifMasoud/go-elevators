package main

import (
	"math"
	"slices"
	"strings"
)

type Elevator struct {
	nextStop        int
	timeLeftOnFloor int
	numFloors       int
	floor           int
	dir             DIRECTION
	upCalls         []bool
	downCalls       []bool
	outCalls        []bool
	doorState       DOORSTATE
}

func NewElevator(nFloors int) *Elevator {
	upCalls := make([]bool, nFloors)
	downCalls := make([]bool, nFloors)
	outCalls := make([]bool, nFloors)
	return &Elevator{numFloors: nFloors, floor: 0, dir: PARKED, upCalls: upCalls, downCalls: downCalls, nextStop: -1, outCalls: outCalls, doorState: CLOSED}
}

func ElevToString(elev Elevator) string {
	var sb strings.Builder

	if elev.doorState == OPEN {
		sb.WriteString("|  |")
	} else {
		sb.WriteString("|--|")
	}
	if elev.dir == UP {
		sb.WriteString("^")
	} else if elev.dir == DOWN {
		sb.WriteString("v")
	}

	return sb.String()
}

func tick(elev *Elevator) {
	if shouldOpen(elev) {
		done := open(elev)
		if done {
			return
		}
	}
	done := waitIfNeeded(elev)
	if done {
		return
	}
	done = closeDoorIfNeeded(elev)
	updateNextStop(elev)
	if done {
		return
	}
	move(elev)
}

func registerCall(elev *Elevator, call Call) {
	switch call.dir {
	case UP:
		elev.upCalls[call.floor] = true
	case DOWN:
		elev.downCalls[call.floor] = true
	case PARKED:
		elev.outCalls[call.floor] = true
	}
}

func move(elev *Elevator) {
	if elev.nextStop == -1 {
		return
	}
	if elev.floor < elev.nextStop {
		elev.floor += 1
	}
	if elev.nextStop < elev.floor {
		elev.floor -= 1
	}
}

func updateNextStop(elev *Elevator) {
	if elev.dir == UP {
		if slices.Contains(aboveUpOutCalls(elev), true) {
			elev.nextStop = slices.IndexFunc(aboveUpOutCalls(elev), func(v bool) bool { return v })
			return
		}
	}
	if elev.dir == DOWN {
		if slices.Contains(belowDownOutCalls(elev), true) {
			elev.nextStop = slices.IndexFunc(belowDownOutCalls(elev), func(v bool) bool { return v })
			return
		}
	}
	elev.nextStop = closestCall(elev)
	matchDirToNextStop(elev)
}

func matchDirToNextStop(elev *Elevator) {
	if elev.nextStop == -1 {
		elev.dir = PARKED
		return
	}
	if elev.upCalls[elev.nextStop] {
		elev.dir = UP
	} else if elev.downCalls[elev.nextStop] {
		elev.dir = DOWN
	} else {
		if elev.floor <= elev.nextStop {
			elev.dir = UP
		} else {
			elev.dir = DOWN
		}
	}

}

func closestCall(elev *Elevator) int {
	closestFlr := -1
	closestDist := math.MaxInt
	megaSlice := elev.outCalls
	for i, b := range elev.downCalls {
		if b && !megaSlice[i] {
			megaSlice[i] = b
		}
	}
	for i, b := range elev.upCalls {
		if b && !megaSlice[i] {
			megaSlice[i] = b
		}
	}
	for f := 0; f < topFloor(elev); f++ {
		if megaSlice[f] {
			if distAandB(elev.floor, f) < closestDist {
				closestFlr = f
				closestDist = distAandB(elev.floor, closestFlr)
			}
		}
	}
	return closestFlr
}

func belowDownOutCalls(elev *Elevator) []bool {
	out := make([]bool, elev.numFloors)
	for i := 0; i <= elev.floor; i++ {
		if elev.outCalls[i] || elev.downCalls[i] {
			out[i] = true
		}
	}
	return out
}

func aboveUpOutCalls(elev *Elevator) []bool {
	out := make([]bool, elev.numFloors)
	for i := elev.floor; i <= topFloor(elev); i++ {
		if elev.outCalls[i] || elev.upCalls[i] {
			out[i] = true
		}
	}
	return out
}

func closeDoorIfNeeded(elev *Elevator) bool {
	if elev.doorState == CLOSED {
		return false
	}
	elev.doorState = CLOSED
	stopDone(elev)
	updateNextStop(elev)
	return true
}

// Returns true if door opened, false if it was already open.
func open(elev *Elevator) bool {
	if elev.doorState == OPEN {
		return false
	}
	elev.timeLeftOnFloor = FLR_WAIT
	elev.doorState = OPEN
	return true
}

func waitIfNeeded(elev *Elevator) bool {
	if elev.timeLeftOnFloor == 0 {
		return false
	}
	elev.timeLeftOnFloor -= 1
	return true
}

func stopDone(elev *Elevator) {
	if elev.nextStop == -1 {
		return
	}
	if elev.dir == UP {
		elev.upCalls[elev.floor] = false
	}
	if elev.dir == DOWN {
		elev.downCalls[elev.floor] = false
	}
	elev.outCalls[elev.floor] = false
	elev.nextStop = -1
}

func shouldOpen(elev *Elevator) bool {
	if elev.dir == UP && elev.nextStop == elev.floor {
		return true
	}
	if elev.dir == DOWN && elev.nextStop == elev.floor {
		return true
	}
	return elev.outCalls[elev.floor]
}

func isPassingByCall(elev *Elevator, call Call) bool {
	return elev.dir == PARKED || elev.dir == UP && call.dir == UP && elev.floor <= call.floor || elev.dir == DOWN && call.dir == DOWN && elev.floor >= call.floor
}

func distToSwitchDir(elev *Elevator) int {
	if elev.dir == UP {
		return topFloor(elev) - elev.floor
	}
	return elev.floor - btmFloor(elev)
}

func topFloor(elev *Elevator) int {
	return elev.numFloors - 1
}

func btmFloor(elev *Elevator) int {
	return 0
}

func endFloor(elev *Elevator) int {
	if elev.dir == UP {
		return topFloor(elev)
	}
	return btmFloor(elev)
}

func oppositeEndFloor(elev *Elevator) int {
	if elev.dir == UP {
		return btmFloor(elev)
	}
	return topFloor(elev)
}

func endToEndDist(elev *Elevator) int {
	return topFloor(elev) - btmFloor(elev)
}

func getDisplacementToCall(elev *Elevator, call Call) int {
	if isPassingByCall(elev, call) {
		return distAandB(elev.floor, call.floor)
	}
	if elev.dir != call.dir {
		return distToSwitchDir(elev) + distAandB(endFloor(elev), call.floor)
	}
	return distToSwitchDir(elev) + endToEndDist(elev) + distAandB(oppositeEndFloor(elev), call.floor)
}

func distAandB(a int, b int) int {
	return int(math.Abs(float64(a) - float64(b)))
}
