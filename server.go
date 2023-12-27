package main

import (
	"math"
	"sync"
)

type Server struct {
	callChan      chan []Call
	elevatorsChan chan []Elevator
	elevators     []*Elevator
	nFloors       int
}

func startServer(server *Server, nFloorsChan chan int, elevsChan chan []Elevator, callsChan chan []Call) {
	server.nFloors = <-nFloorsChan
	resp := SERV_RESP
	nFloorsChan <- resp
	server.callChan = callsChan
	server.elevatorsChan = elevsChan
	for i := 0; i < 3; i++ {
		server.elevators = append(server.elevators, NewElevator(server.nFloors))
	}
	srvLoop(server)
}

func srvLoop(server *Server) {
	// Send initial state:
	server.elevatorsChan <- getElevsValues(server)
	for {
		// Get Calls:
		calls := <-server.callChan
		// Forward to elevators:
		sendCallsElevs(server, calls)
		// Tick:
		tickSrvr(server)
		// Send response:
		server.elevatorsChan <- getElevsValues(server)
	}
}

func tickSrvr(server *Server) {
	var wg sync.WaitGroup
	wg.Add(len(server.elevators))
	for _, elev := range server.elevators {
		elev_ := elev
		go func() {
			tick(elev_)
			defer wg.Done()
		}()
	}
	wg.Wait()
}

func getElevsValues(srvr *Server) []Elevator {
	out := make([]Elevator, 0)
	for _, elev := range srvr.elevators {
		out = append(out, *elev)
	}
	return out
}

func sendCallsElevs(server *Server, calls []Call) {
	for _, call := range calls {
		assignCall(server, call)
	}
}

func assignCall(srv *Server, call Call) {
	if call.elevId != -1 {
		registerCall(srv.elevators[call.elevId], call)
		return
	}
	var bestElev *Elevator = srv.elevators[0]
	minDisplacement := math.MaxInt
	for _, elev := range srv.elevators {
		d := getDisplacementToCall(elev, call)
		if d < minDisplacement || (d == minDisplacement && elev.timeLeftOnFloor < bestElev.timeLeftOnFloor) {
			bestElev = elev
			minDisplacement = d
		}
	}
	registerCall(bestElev, call)
}
