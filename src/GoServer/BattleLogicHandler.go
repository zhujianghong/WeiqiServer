// BattleLogicHandler
package main

import (
	"protobuf"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

func MatchGameHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.MatchGameRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)

}

func InvitePlayGameHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.InvitePlayGameRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	userid := req.GetUserid()
	player := server.playerManager.findPlayerById(userid)
	if player == nil {
		conn.Close()
		return
	}
	log.Println(req)
	res := &protobuf.InvitePlayGameResponse{}
	resMesageid := uint32(protobuf.MsgEnum_InvitePlayGameResponseTag)
	res.Success = protobuf.ErrCodeType_Success.Enum()
	invitePlayGameInfo := req.Info
	messageid := uint32(protobuf.MsgEnum_InvitePlayGameInfoTag)
	recvUserid := invitePlayGameInfo.GetReceiverUserid()
	recvPlayer := server.playerManager.findPlayerById(recvUserid)
	if recvPlayer != nil {
		if recvPlayer.isIdle() {
			recvPlayer.SendMessage(messageid, invitePlayGameInfo)
		} else {
			res.Success = protobuf.ErrCodeType_PlayerPlayingGameError.Enum()
		}
	} else {
		res.Success = protobuf.ErrCodeType_PlayerOfflineError.Enum()
	}
	player.SendMessage(resMesageid, res)
}

func ReceiverPlayGameHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.ReceiverPlayGameRequest{}
	res := &protobuf.ReceiverPlayGameResponse{}
	resMesageid := uint32(protobuf.MsgEnum_ReceiverPlayGameResponseTag)
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	var room *Room = nil

	player := server.playerManager.findPlayerById(req.GetUserid())
	if player == nil {
		conn.Close()
		return
	}
	receiverPlayGameInfo := req.GetInfo()
	messageid := uint32(protobuf.MsgEnum_ReceiverPlayGameInfoTag)
	invitePlayer := server.playerManager.findPlayerById(receiverPlayGameInfo.GetInviterUserid())
	if invitePlayer != nil {
		if req.GetInfo().GetIsAgree() {
			//TODO 建立房间 返回房间id
			if invitePlayer.isPlaying() {
				res.Success = protobuf.ErrCodeType_PlayerPlayingGameError.Enum()
			} else {
				invitePlayer.play()
				player.play()
				room = server.roomManager.createRoom()
				room.initRoom(invitePlayer, player)
				room.beginGame()
				receiverPlayGameInfo.BlackUserid = proto.Int32(room.blackUserid)
				receiverPlayGameInfo.WhiteUserid = proto.Int32(room.whiteUserid)
				receiverPlayGameInfo.RoomId = proto.Uint64(room.roomId)
			}
		} else {
			invitePlayer.SendMessage(messageid, receiverPlayGameInfo)
		}
	} else {
		res.Success = protobuf.ErrCodeType_PlayerOfflineError.Enum()
	}
	player.SendMessage(resMesageid, res)
	if room != nil {
		room.sendMessageToPlayer(0, messageid, receiverPlayGameInfo)
	}
}
