package main

type GoChessStack struct {
	ID           string
	blackWhite   Blackwhite
	lifePoint    int
	chessList    []*GoChess
	gameData     *GameData
	stackManager *StackManager
}

func GetGoChessStack(gameData *GameData, stackManager *StackManager) *GoChessStack {
	var result GoChessStack
	result.gameData = gameData
	result.stackManager = stackManager
	return &result
}

func GetGoChessStackByChess(chess *GoChess, gameData *GameData, stackManager *StackManager) *GoChessStack {
	var result GoChessStack
	result.gameData = gameData
	result.stackManager = stackManager
	result.addChess(chess)
	return &result
}

func GetGoChessStackByChessList(chessList []*GoChess, gameData *GameData, stackManager *StackManager) *GoChessStack {
	var result GoChessStack
	result.gameData = gameData
	result.stackManager = stackManager
	result.addChessList(chessList)
	return &result
}

func (goChessStack *GoChessStack) equal(otherGoChessStack *GoChessStack) bool {
	return goChessStack.ID == otherGoChessStack.ID
}

func (goChessStack *GoChessStack) addChess(goChess *GoChess) {
	goChessStack.chessList = append(goChessStack.chessList, goChess)
	if len(goChessStack.chessList) == 1 {
		goChessStack.ID = goChess.getKey()
	}
}

func (goChessStack *GoChessStack) addChessList(chessList []*GoChess) {
	for _, chess := range chessList {
		goChessStack.addChess(chess)
	}
}

func (goChessStack *GoChessStack) LifePoint() int {
	return goChessStack.lifePoint
}

func (goChessStack *GoChessStack) isChessIn(goChess *GoChess) bool {
	for _, chess := range goChessStack.chessList {
		if chess == goChess {
			return true
		}
	}
	return false
}

func (goChessStack *GoChessStack) megerChessStack(otherGoChessStack *GoChessStack) {
	if goChessStack.blackWhite != otherGoChessStack.blackWhite || goChessStack.ID == otherGoChessStack.ID {
		return
	}
	goChessStack.chessList = append(goChessStack.chessList, otherGoChessStack.chessList...)
}

func (goChessStack *GoChessStack) refreshLifePoint() int {
	return goChessStack.gameData.getStackLifePoint(goChessStack)
}

func (goChessStack *GoChessStack) updateLifePoint() bool {
	return goChessStack.setLifePoint(goChessStack.gameData.getStackLifePoint(goChessStack))
}

func (goChessStack *GoChessStack) setLifePoint(lifePoint int) bool {
	goChessStack.lifePoint = lifePoint
	if lifePoint <= 0 {
		goChessStack.setDeadState()
		return false
	}
	return true
}

func (goChessStack *GoChessStack) lifePointInc() {
	goChessStack.setLifePoint(goChessStack.lifePoint + 1)
}

func (goChessStack *GoChessStack) lifePointdec() {
	goChessStack.setLifePoint(goChessStack.lifePoint - 1)
}

func (goChessStack *GoChessStack) setDeadState() {
	for _, chess := range goChessStack.chessList {
		goChessStack.gameData.removeChess(chess)
		chess.setDeadState()
		chessStack := goChessStack.stackManager.getStackByChess(chess)
		if chessStack != nil {
			goChessStack.stackManager.removeStack(chessStack)
		}
	}
}
