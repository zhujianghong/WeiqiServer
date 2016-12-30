package main

import (
	"math/rand"
	"time"
)

type Robot struct {
	blackWhite      Blackwhite
	gameData        *GameData
	stackManager    *StackManager
	pvpSceneManager *PVPSceneManager
}

func (robot *Robot) putChess() {
	point := robot.getRandPointByBlackwhite(robot.blackWhite)
	robot.gameData.robList = nil
	if point == nil {
		robot.pvpSceneManager.sceneComplete()
		return
	}
	robot.pvpSceneManager.putChess(point)
}

func (robot *Robot) getRandPointByBlackwhite(blackWhite Blackwhite) *GoChessPoint {
	var result *GoChessPoint
	stackList := robot.stackManager.getStacksByBlackWhite(blackWhite)
	stack := getLowestLifePointStack(stackList)
	if stack == nil {
		result = robot.getRandPoint()
	} else {
		result = robot.getRandPointArroundStack(stack)
		var newStacks []*GoChessStack
		if result == nil {
			if len(stackList) == 0 {
				result = robot.getRandPoint()
				return result
			}
			stack = randGetStack(stackList, newStacks)
			if stack == nil {
				result = robot.getRandPoint()
				return result
			}
			result = robot.getRandPointArroundStack(stack)
			//error
			if len(newStacks) <= 0 {
				return result
			}
		}
	}
	return result
}

func (robot *Robot) getRandPointArroundStack(chessStack *GoChessStack) *GoChessPoint {
	var result *GoChessPoint
	if chessStack == nil {
		return nil
	}
	points := robot.GetInvalidPointsArroundStack(chessStack)
	if len(points) > 0 {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		result = points[r.Intn(len(points))]
	}

	if robot.gameData.gameModel == PVE && robot.gameData.currentBlackWhite == robot.blackWhite && result != nil {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if r.Intn(2) == 1 {
			backPoints := robot.GetInvalidPointsArroundPoint(result)
			if len(backPoints) > 0 {
				result = backPoints[r.Intn(len(backPoints))]
			}
		}

	}
	return result
}

func (robot *Robot) getRandPoint() *GoChessPoint {
	points := robot.getAllInvalidPoints()
	if len(points) < 1 {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return points[r.Intn(len(points))]
}

func (robot *Robot) GetInvalidPointsArroundStack(chessStack *GoChessStack) []*GoChessPoint {
	var result []*GoChessPoint
	var flagArray [][][]int = make([][][]int, robot.gameData.boardSize.X)
	for x := 0; x < robot.gameData.boardSize.X; x++ {
		flagArray[x] = make([][]int, robot.gameData.boardSize.Y)
		for y := 0; y < robot.gameData.boardSize.Y; y++ {
			flagArray[x][y] = make([]int, robot.gameData.boardSize.Z)
			for z := 0; z < robot.gameData.boardSize.Z; z++ {
				flagArray[x][y][z] = 0
			}
		}
	}

	for _, chess := range chessStack.chessList {
		points := robot.gameData.getChessSublingList(chess)
		for _, p := range points {
			if !p.hasGoChess() && flagArray[p.X][p.Y][p.Z] == 0 {
				isEnclosur, isOpenEnCloure, isOneLifePoint := robot.gameData.isEnclourePoint(p), robot.gameData.isOpenEnclourePoint(p), robot.gameData.isOnePointLifePoint(p)
				if isEnclosur || isOpenEnCloure || isOneLifePoint {
					continue
				}
				flagArray[p.X][p.Y][p.Z] = 1
				result = append(result, p)
			}
		}
	}
	return result
}

func (robot *Robot) getAllInvalidPoints() []*GoChessPoint {
	var result []*GoChessPoint
	points := robot.gameData.getPointList()
	for _, p := range points {
		if !p.hasGoChess() {
			if robot.gameData.gameModel == PVE && !robot.gameData.isEnclourePoint(p) && !robot.gameData.isOpenEnclourePoint(p) && !robot.gameData.isOnePointLifePoint(p) {
				result = append(result, p)
			}
		}
	}
	return result
}

func (robot *Robot) GetInvalidPointsArroundPoint(point *GoChessPoint) []*GoChessPoint {
	var result []*GoChessPoint
	points := robot.gameData.getPointSublingList(point.X, point.Y, point.Z)
	for _, p := range points {
		if !robot.gameData.isEnclourePoint(p) && !p.hasGoChess() && !robot.gameData.isOpenEnclourePoint(p) && !robot.gameData.isOnePointLifePoint(p) {
			result = append(result, p)
		}
	}
	return result
}
