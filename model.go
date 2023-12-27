package main

const SERV_RESP = 88
const FLR_WAIT = 3

type Call struct {
	floor  int
	dir    DIRECTION
	elevId int
}

type DIRECTION int

const (
	UP     DIRECTION = iota
	DOWN   DIRECTION = iota
	PARKED DIRECTION = iota
)

type DOORSTATE int

const (
	CLOSED DOORSTATE = iota
	OPEN   DOORSTATE = iota
)
