package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	nFloors   int
	elevsChan chan []Elevator
	callsChan chan []Call
}

func NewClient(nFloors int) *Client {
	return &Client{nFloors: nFloors, elevsChan: make(chan []Elevator), callsChan: make(chan []Call)}
}

func sendEvent(client *Client, calls []Call) {
	client.callsChan <- calls
}

func run(clnt *Client) {
	initConnectionToServer(clnt)
	startLoop(clnt)
}

func startLoop(clnt *Client) {
	// Get initial state:
	elevs := <-clnt.elevsChan
	writeElevatorsFiles(elevs)
	for {
		// Get calls from input:
		var calls []Call = getCallsFromInput()
		// Send Calls:
		clnt.callsChan <- calls
		// Get Response:
		elevs := <-clnt.elevsChan
		// Render response:
		writeElevatorsFiles(elevs)

	}
}

func getCallsFromInput() []Call {
	calls := make([]Call, 0)
	for {
		fmt.Println("Enter next call (format={floorNum} {U/D/P/elevId}, leave empty when done:")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		if input == "\n" {
			break
		}
		var err error = nil
		callFlds := strings.Fields(input)
		fmt.Println("Read: floor="+callFlds[0]+", Dir/Elev=", callFlds[1])
		flr, err := strconv.ParseInt(callFlds[0], 10, 32)
		check(err)
		var dir DIRECTION
		var elevId int64 = -1
		if strings.ToLower(callFlds[1]) == "u" {
			dir = UP
		} else if strings.ToLower(callFlds[1]) == "d" {
			dir = DOWN
		} else {
			dir = PARKED
			elevId, err = strconv.ParseInt(callFlds[1], 10, 32)
		}
		call := Call{
			floor:  int(flr),
			dir:    dir,
			elevId: int(elevId),
		}
		calls = append(calls, call)
	}
	return calls
}

func initConnectionToServer(clnt *Client) error {
	c := make(chan int)
	srvr := Server{}
	go startServer(&srvr, c, clnt.elevsChan, clnt.callsChan)
	c <- clnt.nFloors
	resp := <-c
	if resp != SERV_RESP {
		return errors.New("Invalid Server Connection response")
	}
	return nil
}
