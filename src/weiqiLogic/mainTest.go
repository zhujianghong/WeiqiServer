package main

import (
	"fmt"
	"log"
)

func main() {
	pvpSceneManager := getPVPSceneManagerByModel(PVE)
	fmt.Println("Game begin, please input coordinate like : x, y, z.")
	pvpSceneManager.startNewGame()
	var x, y, z int
	for {
		_, ok := fmt.Scanf("%d%d%d\n", &x, &y, &z)
		if ok != nil {
			log.Println(ok)
			break
		}
		point := pvpSceneManager.findPointByPos(x, y, z)
		pvpSceneManager.putChess(point)
		fmt.Println(pvpSceneManager.gameData.getEmptyPointNum())
	}
	fmt.Println("abc")
}
