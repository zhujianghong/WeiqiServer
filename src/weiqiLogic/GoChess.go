package main

import (
	"fmt"
)

type Blackwhite int

const (
	BLACK = 0
	WHITE = 1
)

type GoChess struct {
	blackWhite   Blackwhite
	gochessPoint *GoChessPoint
	X, Y, Z      int
}

func (gochess *GoChess) init(pos Vector3D, blackWhite Blackwhite) {
	gochess.X = int(pos.X)
	gochess.Y = int(pos.Y)
	gochess.Z = int(pos.Z)
	gochess.blackWhite = blackWhite
}

func (gochess *GoChess) getKey() string {
	return fmt.Sprintf("%d,%d,%d", gochess.X, gochess.Y, gochess.Z)
}

func (gochess *GoChess) setDeadState() {
	gochess.gochessPoint.Gochess = nil
	gochess.gochessPoint = nil
}

var allowTouch bool = true

type GoChessPoint struct {
	X, Y, Z int
	Gochess *GoChess
}

func (point *GoChessPoint) init(pos Vector3D) {
	point.X = int(pos.X)
	point.Y = int(pos.Y)
	point.Z = int(pos.Z)
}

func (point *GoChessPoint) hasGoChess() bool {
	if point.Gochess == nil {
		return false
	}
	return true
}
