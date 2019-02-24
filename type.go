package main

const Dim = 25

var Special bool = false
var Points int

type Position struct {
	X, Y byte
}

type Player struct {
	Alive bool
	Dir   byte
	Pos   Position
}

type Ghost struct {
	Alive bool
	Dir   byte
	Pos   Position
}
type Lab [Dim][Dim]int
