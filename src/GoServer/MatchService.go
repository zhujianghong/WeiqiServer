// MatchService
package main

import (
	"time"
)

type MatchPlayerInfo struct {
	userid int32
	grade  int32
}

type MatchService struct {
	MatchQueue chan *MatchPlayerInfo
	stopFlag   chan struct{}
}

func newMatchService() *MatchService {
	return &MatchService{
		MatchQueue: make(chan *MatchPlayerInfo, 500),
		stopFlag:   make(chan struct{}),
	}
}

func (m *MatchService) Start() {
	matchingPlayerMap := make(map[int32]*MatchPlayerInfo)
	for {
		select {
		case <-m.stopFlag:
			return

		case <-time.After(1000 * time.Millisecond):

		case info := <-m.MatchQueue:
			if _, ok := matchingPlayerMap[info.userid]; !ok {
				matchingPlayerMap[info.userid] = info
			}
		}
		//开始匹配
		length := len(matchingPlayerMap)
		useridList := make([]int32, length)
		for k, _ := range matchingPlayerMap {
			length--
			useridList[length] = k
			if length == 0 {
				break
			}
		}
		for i := 0; i < length; i = i + 2 {
			player := server.playerManager.findPlayerById(useridList[i])
			if player != nil && i+1 != length {
				otherPlayer := server.playerManager.findPlayerById(useridList[i+1])
				if otherPlayer != nil {
					room := server.roomManager.createRoom()
					room.initRoom(player, otherPlayer)
					room.beginGame()
					room.SendMatchSuccessInfo()
					//TODO 通知两个玩家的信息和房间id
					delete(matchingPlayerMap, useridList[i])
					delete(matchingPlayerMap, useridList[i+1])
				}
			}
		}
	}
}

func (m *MatchService) Stop() {
	close(m.stopFlag)
}
