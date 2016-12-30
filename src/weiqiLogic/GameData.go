package main

const (
	DefaultChessBoard = 11
)

type GameModel int // 0-pvp, 1-pve
type GameRole int  // 0-human, 1-robot

const (
	PVP   = 0
	PVE   = 1
	HUMAN = 0
	ROBOT = 1
)

type GameData struct {
	boardSize         ChessBoardSize
	stackManager      *StackManager
	currentBlackWhite Blackwhite // 0-黑， 1-白
	gameModel         GameModel
	whiteChesses      map[string]*GoChess
	blackChesses      map[string]*GoChess
	chessBoardPoints  [][][]*GoChessPoint
	pointList         []*GoChessPoint
	robList           []*GoChessPoint
}

func GetGameData(boardSize *ChessBoardSize, model GameModel) *GameData {
	var result GameData
	result.gameModel = model
	if boardSize == nil {
		result.boardSize = GetChessBoardSize(DefaultChessBoard)
	}
	result.chessBoardPoints = make([][][]*GoChessPoint, result.boardSize.X)
	for x := 0; x < result.boardSize.X; x++ {
		result.chessBoardPoints[x] = make([][]*GoChessPoint, result.boardSize.Y)
		for y := 0; y < result.boardSize.Y; y++ {
			result.chessBoardPoints[x][y] = make([]*GoChessPoint, result.boardSize.Z)
		}
	}
	result.reset()
	return &result
}

func (gameData *GameData) reset() {
	gameData.currentBlackWhite = 0
	gameData.whiteChesses = nil
	gameData.blackChesses = nil
	gameData.robList = nil
}

func (gameData *GameData) getPointList() []*GoChessPoint {
	if gameData.pointList != nil {
		return gameData.pointList
	}
	for _, x := range gameData.chessBoardPoints {
		for _, xy := range x {
			for _, xyz := range xy {
				if xyz != nil {
					gameData.pointList = append(gameData.pointList, xyz)
				}
			}
		}
	}
	return gameData.pointList
}

func (gameData *GameData) matchSettle() Vector2D {
	var blackEyes, whiteEyes, allPoints, emptyPoints []*GoChessPoint
	allPoints = gameData.getPointList()
	for _, point := range allPoints {
		if !point.hasGoChess() {
			emptyPoints = append(emptyPoints, point)
		}
	}
	var samePoints, diffPoints []*GoChessPoint
	for _, point := range emptyPoints {
		gameData.getPointSibling(point, &samePoints, &diffPoints)
		if len(samePoints) > 0 && len(samePoints) > len(diffPoints) {
			if samePoints[0].hasGoChess() {
				if samePoints[0].Gochess.blackWhite == 0 {
					blackEyes = append(blackEyes, point)
				} else {
					whiteEyes = append(whiteEyes, point)
				}

			}
		} else if len(diffPoints) > 0 && len(diffPoints) > len(samePoints) {
			if diffPoints[0].hasGoChess() {
				if diffPoints[0].Gochess.blackWhite == 0 {
					blackEyes = append(blackEyes, point)
				} else {
					whiteEyes = append(whiteEyes, point)
				}
			}
		}
	}

	return Vector2D{
		float32(len(blackEyes) + len(gameData.blackChesses)),
		float32(len(whiteEyes) + len(gameData.whiteChesses)),
	}
}

func (gameData *GameData) getChessNum() (int, int) {
	return len(gameData.blackChesses), len(gameData.whiteChesses)
}

func (gameData *GameData) getStackLifePoint(chessStack *GoChessStack) int {
	return len(gameData.getStackLifePointList(chessStack))
}

func (gameData *GameData) getStackLifePointList(chessStack *GoChessStack) []*GoChessPoint {
	var points []*GoChessPoint
	flag := make([][][]int, gameData.boardSize.X)
	for x := 0; x < gameData.boardSize.X; x++ {
		flag[x] = make([][]int, gameData.boardSize.Y)
		for y := 0; y < gameData.boardSize.Y; y++ {
			flag[x][y] = make([]int, gameData.boardSize.Z)
		}
	}
	for _, chess := range chessStack.chessList {
		pointList := gameData.getChessSublingList(chess)
		for _, point := range pointList {
			if !point.hasGoChess() && flag[point.X][point.Y][point.Z] == 0 {
				flag[point.X][point.Y][point.Z] = 1
				points = append(points, point)
			}
		}
	}
	return points
}
func (gameData *GameData) getPointByChess(chess *GoChess) *GoChessPoint {
	point := chess.gochessPoint
	if point == nil {
		point = gameData.getPointByPos(chess.X, chess.Y, chess.Z)
	}
	return point
}

func (gameData *GameData) getPointByPos(x, y, z int) *GoChessPoint {
	points := gameData.getPointList()
	for _, point := range points {
		if point.X == x && point.Y == y && point.Z == z {
			return point
		}
	}
	return nil
}

func (gameData *GameData) isEnclourePoint(point *GoChessPoint) bool {
	var samePoints, diffPoints, pointList []*GoChessPoint
	gameData.getPointSibling(point, &samePoints, &diffPoints)
	pointList = gameData.getPointSublingList(point.X, point.Y, point.Z)
	if len(diffPoints) == len(pointList) {
		return true
	}
	return false
}

func (gameData *GameData) isOpenEnclourePoint(point *GoChessPoint) bool {
	var samePoints, diffPoints, pointList []*GoChessPoint
	gameData.getPointSibling(point, &samePoints, &diffPoints)
	pointList = gameData.getPointSublingList(point.X, point.Y, point.Z)
	if len(samePoints) == len(pointList) {
		return true
	}
	return false
}

func (gameData *GameData) isRobPoint(point *GoChessPoint) bool {
	for _, p := range gameData.robList {
		if p == point {
			return true
		}
	}
	return false
}

func (gameData *GameData) isOnePointLifePoint(point *GoChessPoint) bool {
	var samePoints, diffPoints []*GoChessPoint
	gameData.getPointSibling(point, &samePoints, &diffPoints)
	if len(samePoints) < 1 {
		return false
	}
	//todo
	return false
}

func (gameData *GameData) isOnePointLifeStack(stack *GoChessStack) bool {
	if stack == nil {
		return false
	}
	if stack.LifePoint() <= 1 {
		return true
	}
	return false
}

func (gameData *GameData) getPointLifePoint(point *GoChessPoint) int {
	points := gameData.getPointSublingList(point.X, point.Y, point.Z)
	var lifePoint int
	for _, p := range points {
		if p != nil && p.hasGoChess() {
			lifePoint++
		}
	}
	return lifePoint
}

func (gameData *GameData) checkPublicArea() bool {
	points := gameData.getPointList()
	for _, point := range points {
		if point.hasGoChess() {
			continue
		}
		var samePoints, diffPoints []*GoChessPoint
		gameData.getPointSibling(point, &samePoints, &diffPoints)
		if len(samePoints) < 4 || len(diffPoints) < 4 {
			return true
		}
	}
	return false
}

func (gameData *GameData) getPointSibling(point *GoChessPoint, sPoints, dPoints *[]*GoChessPoint) {
	if point == nil {
		return
	}

	if point.Gochess != nil {
		chess := point.Gochess
		gameData.getChessSibling(chess, sPoints, dPoints)
		return
	} else {
		points := gameData.getPointSublingList(point.X, point.Y, point.Z)
		for _, p := range points {
			if p.hasGoChess() {
				if p.Gochess.blackWhite == gameData.currentBlackWhite {
					*sPoints = append(*sPoints, p)
				} else {
					*dPoints = append(*dPoints, p)
				}
			}
		}
	}
}

func (gameData *GameData) getChessSibling(chess *GoChess, sPoints, dPoints *[]*GoChessPoint) {
	points := gameData.getChessSublingList(chess)
	for _, point := range points {
		if point.hasGoChess() {
			if point.Gochess.blackWhite == chess.blackWhite {
				*sPoints = append(*sPoints, point)
			} else {
				*dPoints = append(*dPoints, point)
			}
		}
	}
}

func (gameData *GameData) getChessSublingList(chess *GoChess) []*GoChessPoint {
	var result []*GoChessPoint
	result = gameData.getPointSublingList(chess.X, chess.Y, chess.Z)
	return result
}

func (gameData *GameData) getPointSublingList(x, y, z int) []*GoChessPoint {
	var result []*GoChessPoint
	if gameData.preCheck(x-1, y, z) && gameData.chessBoardPoints[x-1][y][z] != nil {
		result = append(result, gameData.chessBoardPoints[x-1][y][z])
	}
	if gameData.preCheck(x+1, y, z) && gameData.chessBoardPoints[x+1][y][z] != nil {
		result = append(result, gameData.chessBoardPoints[x+1][y][z])
	}
	if gameData.preCheck(x, y-1, z) && gameData.chessBoardPoints[x][y-1][z] != nil {
		result = append(result, gameData.chessBoardPoints[x][y-1][z])
	}
	if gameData.preCheck(x, y+1, z) && gameData.chessBoardPoints[x][y+1][z] != nil {
		result = append(result, gameData.chessBoardPoints[x][y+1][z])
	}
	if gameData.preCheck(x, y, z-1) && gameData.chessBoardPoints[x][y][z-1] != nil {
		result = append(result, gameData.chessBoardPoints[x][y][z-1])
	}
	if gameData.preCheck(x, y, z+1) && gameData.chessBoardPoints[x][y][z+1] != nil {
		result = append(result, gameData.chessBoardPoints[x][y][z+1])
	}
	return result
}

func (gameData *GameData) preCheck(x, y, z int) bool {
	if x < 0 || x > gameData.boardSize.X-1 || y < 0 || y > gameData.boardSize.Y-1 || z < 0 || z > gameData.boardSize.Z-1 {
		return false
	}
	return true
}

func (gameData *GameData) changeBlackWhite() Blackwhite {
	if gameData.currentBlackWhite == BLACK {
		gameData.currentBlackWhite = WHITE
	} else {
		gameData.currentBlackWhite = BLACK
	}
	return gameData.currentBlackWhite
}

func (gameData *GameData) addChess(chess *GoChess) {
	if gameData.whiteChesses == nil {
		gameData.whiteChesses = make(map[string]*GoChess)
	}
	if gameData.blackChesses == nil {
		gameData.blackChesses = make(map[string]*GoChess)
	}

	_, bOk := gameData.whiteChesses[chess.getKey()]
	_, wOk := gameData.blackChesses[chess.getKey()]
	if !bOk && !wOk {
		switch chess.blackWhite {
		case 0:
			gameData.blackChesses[chess.getKey()] = chess
		case 1:
			gameData.whiteChesses[chess.getKey()] = chess
		}
	}
}

func (gameData *GameData) removeChess(chess *GoChess) {
	switch chess.blackWhite {
	case 0:
		delete(gameData.blackChesses, chess.getKey())
	case 1:
		delete(gameData.whiteChesses, chess.getKey())
	}
}

func (gameData *GameData) getEmptyPointNum() int {
	var num int
	for _, point := range gameData.pointList {
		if !point.hasGoChess() {
			num++
		}
	}
	return num
}

type ChessBoardSize struct {
	X, Y, Z    int
	pointCount int
}

func GetChessBoardSize(values ...int) ChessBoardSize {
	if len(values) != 1 && len(values) != 3 {
		panic("wrong args number.")
	}
	var result ChessBoardSize
	if len(values) == 1 {
		result.X = values[0]
		result.Y = values[0]
		result.Z = values[0]
	}
	if len(values) == 3 {
		result.X = values[0]
		result.Y = values[1]
		result.Z = values[2]
	}
	result.pointCount = result.X*result.Y*result.Z - ((result.X - 2) * (result.Y - 2) * (result.Z - 2))
	return result
}
