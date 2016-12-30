package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Player struct {
	blackWhite Blackwhite
}

type PVPSceneManager struct {
	gameData           GameData
	stackManager       StackManager
	player1            Player
	player2            Player
	robot              Robot
	gameModel          GameModel
	lastPoint          *GoChessPoint
	localLastViewPoint *GoChessPoint
	isPosLabInShow     bool
}

func getPVPSceneManagerByModel(model GameModel) *PVPSceneManager {
	var result = PVPSceneManager{
		gameData:           *GetGameData(nil, model),
		stackManager:       StackManager{nil, nil, nil},
		player1:            Player{},
		gameModel:          model,
		lastPoint:          nil,
		localLastViewPoint: nil,
		isPosLabInShow:     false,
	}
	result.stackManager.gameData = &result.gameData
	result.gameData.stackManager = &result.stackManager
	if model == PVP {
		result.player2 = Player{}
	} else {
		result.robot.gameData = &result.gameData
		result.robot.stackManager = &result.stackManager
		result.robot.pvpSceneManager = &result
	}
	return &result
}

func (pvpSceneManager *PVPSceneManager) startNewGame() {
	pvpSceneManager.gameData.reset()
	pvpSceneManager.stackManager.reset()
	pvpSceneManager.lastPoint = nil

	temp := pvpSceneManager.gameData.getPointList()
	if len(temp) > 1 {
		pvpSceneManager.setAllowTouch(true)
	} else {
		pvpSceneManager.createChessBoard()
	}
	pvpSceneManager.randBlackWhite()
}

func (pvpSceneManager *PVPSceneManager) putChess(point *GoChessPoint) {
	if point.hasGoChess() {
		//todo
		return
	}

	isRob, isEnclosure, catEat, willKillHimself := pvpSceneManager.gameData.isRobPoint(point), pvpSceneManager.gameData.isEnclourePoint(point), pvpSceneManager.stackManager.canEatChessPoint(point), pvpSceneManager.stackManager.willKillHimself(point)
	if isRob || (isEnclosure && !catEat) || (willKillHimself && !catEat) {
		//todo
		return
	}

	chess := GoChess{}
	chess.init(Vector3D{float32(point.X), float32(point.Y), float32(point.Z)}, pvpSceneManager.gameData.currentBlackWhite)
	chess.gochessPoint = point

	point.Gochess = &chess
	if chess.blackWhite == BLACK {
		fmt.Println("Black put %d,%d,%d", chess.X, chess.Y, chess.Z)
	}
	if chess.blackWhite == WHITE {
		fmt.Println("white put %d,%d,%d", chess.X, chess.Y, chess.Z)
	}

	if !point.hasGoChess() {
		panic("")
	}

	pvpSceneManager.lastPoint = point
	stacksToclear := pvpSceneManager.stackManager.placeChess(&chess)
	stacksToclear = stacksToclear
	hasPublicArea := pvpSceneManager.gameData.checkPublicArea()
	if !hasPublicArea {
		pvpSceneManager.sceneComplete()
	}

	nextBlackWhite := pvpSceneManager.gameData.changeBlackWhite()
	if nextBlackWhite == pvpSceneManager.player1.blackWhite {
		//轮到1下棋  发送轮着信息给客户端
	} else {
		//轮到2下棋  发送轮着信息给客户端
		if pvpSceneManager.gameModel == PVE {
			pvpSceneManager.robot.putChess()
		}
	}
}

func (pvpSceneManager *PVPSceneManager) randBlackWhite() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	firstHandPlayer := r.Intn(2)
	if firstHandPlayer == 0 {
		pvpSceneManager.player1.blackWhite = BLACK
	} else {
		pvpSceneManager.player1.blackWhite = WHITE
	}
	if pvpSceneManager.gameModel == PVP {
		pvpSceneManager.player2.blackWhite = 1 - pvpSceneManager.player1.blackWhite
	} else {
		pvpSceneManager.robot.blackWhite = 1 - pvpSceneManager.player1.blackWhite
		if pvpSceneManager.robot.blackWhite == BLACK {
			pvpSceneManager.robot.putChess()
		}

	}
}

func (pvpSceneManager *PVPSceneManager) setBoardSize(size int) {
	pvpSceneManager.gameData.boardSize = GetChessBoardSize(size)
}

func (pvpSceneManager *PVPSceneManager) sceneComplete() {
	pvpSceneManager.onSceneComplete()
}

func (pvpSceneManager *PVPSceneManager) createChessBoard() {
	chessBoardSize := pvpSceneManager.gameData.boardSize
	pvpSceneManager.gameData.chessBoardPoints = make([][][]*GoChessPoint, chessBoardSize.X)
	for x := 0; x < chessBoardSize.X; x++ {
		pvpSceneManager.gameData.chessBoardPoints[x] = make([][]*GoChessPoint, chessBoardSize.Y)
		for y := 0; y < chessBoardSize.Y; y++ {
			pvpSceneManager.gameData.chessBoardPoints[x][y] = make([]*GoChessPoint, chessBoardSize.Z)
		}
	}

	for x := 0; x < chessBoardSize.X; x++ {
		for y := 0; y < chessBoardSize.Y; y++ {
			for z := 0; z < chessBoardSize.Z; z++ {
				if x == 0 || y == 0 || z == 0 || x == chessBoardSize.X-1 || y == chessBoardSize.Y-1 || z == chessBoardSize.Z-1 {
					point := GoChessPoint{x, y, z, nil}
					pvpSceneManager.gameData.chessBoardPoints[x][y][z] = &point
				}
			}
		}
	}
	pvpSceneManager.gameData.getPointList()
	pvpSceneManager.setAllowTouch(true)
}

func (pvpSceneManager *PVPSceneManager) getPlayerBlackWhite() Blackwhite {
	return pvpSceneManager.player1.blackWhite
}

func (pvpSceneManager *PVPSceneManager) setAllowTouch(allow bool) {
	allowTouch = allow
}

func (pvpSceneManager *PVPSceneManager) findPointByPos(x, y, z int) *GoChessPoint {
	var result *GoChessPoint = nil
	pointList := pvpSceneManager.gameData.getPointList()
	for _, point := range pointList {
		if point.X == x && point.Y == y && point.Z == z {
			result = point
			break
		}
	}
	return result
}

func (pvpSceneManager *PVPSceneManager) getInputLocalToReal(x, y, z int) Vector3D {
	offset := int(pvpSceneManager.gameData.boardSize.X / 2)
	result := Vector3D{
		float32(x + offset),
		float32(y + offset),
		float32(z + offset),
	}
	return result
}

func (pvpSceneManager *PVPSceneManager) putChessByLocalInput(x, y, z int) {
	pos := Vector3D{float32(x), float32(y), float32(z)}
	if vector3DEqual(pos, scale3D(One3D, float32(pvpSceneManager.gameData.boardSize.X))) {
		return
	}
	point := pvpSceneManager.findPointByPos(x, y, z)
	if point == nil || point.hasGoChess() {
		return
	}
	pvpSceneManager.putChess(point)
}

func (pvpSceneManager *PVPSceneManager) onSceneComplete() {
	var result int
	blackEyes, whiteEyes := pvpSceneManager.gameData.getChessNum()
	if blackEyes == whiteEyes {
		result = 0
	} else {
		if blackEyes > whiteEyes && pvpSceneManager.player1.blackWhite == 0 || blackEyes < whiteEyes && pvpSceneManager.player1.blackWhite == 1 {
			result = 1
		} else {
			result = 0
		}
	}
	result = result
}
