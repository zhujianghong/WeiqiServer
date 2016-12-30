package main

import (
	"math/rand"
	"time"
)

type StackManager struct {
	gameData                                 *GameData
	blackChessStackList, whiteChessStackList []*GoChessStack
}

func (stackManager *StackManager) reset() {
	stackManager.blackChessStackList, stackManager.whiteChessStackList = nil, nil
}

func (stackManager *StackManager) addStack(chessStack *GoChessStack) {
	switch chessStack.blackWhite {
	case 0:
		stackManager.blackChessStackList = append(stackManager.blackChessStackList, chessStack)
	case 1:
		stackManager.whiteChessStackList = append(stackManager.whiteChessStackList, chessStack)
	}
}

func (stackManager *StackManager) removeStack(chessStack *GoChessStack) {
	if chessStack == nil {
		return
	}
	switch chessStack.blackWhite {
	case 0:
		index := findStackIndex(stackManager.blackChessStackList, chessStack)
		if index != -1 {
			stackManager.blackChessStackList = append(stackManager.blackChessStackList[:index], stackManager.blackChessStackList[index+1:]...)
		}
	case 1:
		index := findStackIndex(stackManager.whiteChessStackList, chessStack)
		if index != -1 {
			stackManager.whiteChessStackList = append(stackManager.whiteChessStackList[:index], stackManager.whiteChessStackList[index+1:]...)
		}
	}
}

func (stackManager *StackManager) getStackByChess(chess *GoChess) *GoChessStack {
	var result *GoChessStack
	switch chess.blackWhite {
	case 0:
		result = getStackByChess(stackManager.blackChessStackList, chess)
	case 1:
		result = getStackByChess(stackManager.whiteChessStackList, chess)
	}
	return result
}

func (stackManager *StackManager) mergeTwoStack(chessStack1, chessStack2 *GoChessStack) {
	if chessStack1 == nil || chessStack2 == nil {
		return
	}
	chessStack1.megerChessStack(chessStack2)
	stackManager.removeStack(chessStack2)
}

func (stackManager *StackManager) mergeStack(stackList []*GoChessStack) *GoChessStack {
	if len(stackList) == 0 {
		return nil
	}
	chessStackDst := stackList[0]
	for _, stack := range stackList[1:] {
		stackManager.mergeTwoStack(chessStackDst, stack)
	}
	return chessStackDst
}

func (stackManager *StackManager) splitStack(sourceStack *GoChessStack, splitChess *GoChess) []*GoChessStack {
	var result []*GoChessStack
	if !sourceStack.isChessIn(splitChess) {
		return result
	}
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getChessSibling(splitChess, &samePoints, &diffPoints)
	for _, point := range samePoints {
		var temp []*GoChess
		temp = append(temp, point.Gochess)
		temp = stackManager.chessGrow(&temp, point.Gochess)
		pTemp := GetGoChessStack(stackManager.gameData, stackManager)
		pTemp.addChessList(temp)
		result = append(result, pTemp)
	}
	stackManager.removeStack(sourceStack)
	for _, stack := range result {
		stackManager.addStack(stack)
	}
	stackManager.gameData.removeChess(splitChess)
	return result
}

func (stackManager *StackManager) chessGrow(chessList *[]*GoChess, chess *GoChess) []*GoChess {
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getChessSibling(chess, &samePoints, &diffPoints)
	for _, point := range samePoints {
		index := findChessIndex(*chessList, point.Gochess)
		if index == -1 {
			*chessList = append(*chessList, point.Gochess)
		}
	}
	return *chessList
}

func (stackManager *StackManager) getStacksByPoints(pointList []*GoChessPoint) []*GoChessStack {
	var resultMap map[*GoChessStack]bool = make(map[*GoChessStack]bool)
	var result []*GoChessStack
	for _, point := range pointList {
		resultMap[stackManager.getStackByChess(point.Gochess)] = true
	}
	for k := range resultMap {
		if k == nil {
			continue
		}
		result = append(result, k)
	}
	return result
}

func (stackManager *StackManager) getStacksByBlackWhite(blackWhite Blackwhite) []*GoChessStack {
	var result []*GoChessStack
	switch blackWhite {
	case 0:
		result = stackManager.blackChessStackList
	case 1:
		result = stackManager.whiteChessStackList
	}
	return result
}

func (stackManager *StackManager) getDiffStacksByChess(chess *GoChess) []*GoChessStack {
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getChessSibling(chess, &samePoints, &diffPoints)
	return stackManager.getStacksByPoints(diffPoints)
}

func (stackManager *StackManager) placeChess(chess *GoChess) []*GoChessStack {
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getChessSibling(chess, &samePoints, &diffPoints)
	var chessStacks, chessStacksToClear []*GoChessStack
	var chessStack *GoChessStack
	//var hadEat bool = false
	if len(samePoints) == 0 {
		chessStack = GetGoChessStackByChess(chess, stackManager.gameData, stackManager)
		stackManager.addStack(chessStack)
		chessStack.refreshLifePoint()

		chessStacks = stackManager.getDiffStacksByChess(chess)
		point := chess.gochessPoint
		canEatPoint := stackManager.canEatChessPoint(point)
		if !canEatPoint && !chessStack.updateLifePoint() {
			stackManager.clearSite(chessStack)
			stackManager.gameData.removeChess(chess)
			return nil
		}
	} else {
		var stack *GoChessStack
		chessStack = GetGoChessStackByChess(chess, stackManager.gameData, stackManager)
		chessStacks = stackManager.getStacksByPoints(samePoints)
		chessStacks = append(chessStacks, chessStack)
		stack = stackManager.mergeStack(chessStacks)
		stack.updateLifePoint()
	}
	stackManager.gameData.addChess(chess)
	chessStacks = stackManager.getStacksByPoints(diffPoints)
	for _, stack := range chessStacks {
		stack.lifePointdec()
		if stack.LifePoint() <= 0 {
			stackManager.gameData.robList = nil
			if len(stack.chessList) == 1 {
				point := stackManager.gameData.getPointByChess(stack.chessList[0])
				stackManager.gameData.robList = append(stackManager.gameData.robList, point)
			}
			chessStacksToClear = append(chessStacksToClear, stack)
			stackManager.clearSite(stack)
		}
	}
	return chessStacksToClear
}

func (stackManager *StackManager) canEatChessPoint(point *GoChessPoint) bool {
	var result bool = false
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getPointSibling(point, &samePoints, &diffPoints)
	chessStacks := stackManager.getStacksByPoints(diffPoints)
	for _, stack := range chessStacks {
		if stackManager.gameData.getStackLifePoint(stack) <= 1 {
			result = true
		}
	}
	return result
}

func (stackManager *StackManager) willKillHimself(point *GoChessPoint) bool {
	var result bool
	var samePoints, diffPoints []*GoChessPoint
	stackManager.gameData.getPointSibling(point, &samePoints, &diffPoints)
	stackList := stackManager.getStacksByPoints(samePoints)
	if len(stackList) < 1 || stackManager.gameData.getPointLifePoint(point) > 0 {
		return false
	}
	stackList = stackManager.getStacksByPoints(samePoints)
	for _, stack := range stackList {
		lifePoint := stackManager.gameData.getStackLifePoint(stack)
		if lifePoint <= 1 {
			result = true
		} else {
			result = false
		}
	}
	return result
}

func (stackManager *StackManager) clearSite(chessStack *GoChessStack) {
	stackManager.addStackLifePoint(chessStack)
	stackManager.removeStack(chessStack)
}

func (stackManager *StackManager) addStackLifePoint(chessStack *GoChessStack) {
	for _, chess := range chessStack.chessList {
		stackList := stackManager.getDiffStacksByChess(chess)
		for _, stack := range stackList {
			if stack != nil {
				stack.lifePointInc()
			}
		}
	}
}

func getStackByChess(stackList []*GoChessStack, chess *GoChess) *GoChessStack {
	for _, chessStack := range stackList {
		if chessStack.isChessIn(chess) {
			return chessStack
		}
	}
	return nil
}

func findStackIndex(stackList []*GoChessStack, chessStack *GoChessStack) int {
	for index, stack := range stackList {
		if stack == chessStack {
			return index
		}
	}
	return -1
}

func findChessIndex(chessList []*GoChess, goChess *GoChess) int {
	for index, chess := range chessList {
		if chess == goChess {
			return index
		}
	}
	return -1
}

func getLowestLifePointStack(stackList []*GoChessStack) *GoChessStack {
	if stackList == nil || len(stackList) == 0 {
		return nil
	}
	var result *GoChessStack
	for idx, stack := range stackList[1:] {
		if idx == 0 {
			result = stack
		} else {
			if result.LifePoint() > stack.LifePoint() {
				result = stack
			}
		}
	}
	return result
}

func ranGetStack(stackList []*GoChessStack) *GoChessStack {
	if stackList == nil || len(stackList) == 0 {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return stackList[r.Intn(len(stackList))]
}

func randGetStack(stackList, newStackList []*GoChessStack) *GoChessStack {
	if stackList == nil || len(stackList) == 0 {
		return nil
	}
	lowestLifePointStack := getLowestLifePointStack(stackList)
	if newStackList == nil {
		newStackList = stackList
	}
	index := findStackIndex(newStackList, lowestLifePointStack)
	if index != -1 {
		newStackList = append(newStackList[:index], newStackList[index+1:]...)
	}
	if len(newStackList) <= 1 {
		return nil
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index = r.Intn(len(newStackList))
	var result *GoChessStack
	if newStackList[index] != nil {
		result = newStackList[index]
		newStackList = append(newStackList[:index], newStackList[index+1:]...)
	}
	return result
}
