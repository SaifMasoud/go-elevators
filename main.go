package main

import (
	"os"
	"strconv"
)

func main() {
	clnt := NewClient(10)
	run(clnt)
}

func writeElevatorsFiles(elevators []Elevator) {
	var fnames []string
	for i := 0; i < len(elevators); i++ {
		fnames = append(fnames, "e"+strconv.FormatInt(int64(i), 10)+".txt")
		writeElevatorToFile(elevators[i], fnames[i])
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func writeElevatorToFile(elevator Elevator, fname string) {
	f, err := os.Create(fname)
	check(err)
	defer f.Close()

	for i := elevator.numFloors - 1; i >= 0; i-- {
		f.WriteString(strconv.FormatInt(int64(i), 10) + ": ")
		if i == elevator.floor {
			f.WriteString(ElevToString(elevator))
		}
		f.WriteString("\n")
	}
}
