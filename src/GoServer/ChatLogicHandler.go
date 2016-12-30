// ChatLogicHandler
package main

import (
	"protobuf"
	"time"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

func ChatHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	userid := conn.GetExtraData().(*PlayerData).GetUserid()
	player := server.playerManager.findPlayerById(userid)
	if player == nil {
		conn.Close()
		return
	}
	req := &protobuf.ChatRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	res := &protobuf.ChatResponse{}
	res.Success = protobuf.ErrCodeType_Success.Enum()
	chatMessage := req.Info
	ChatMessageid := uint32(protobuf.MsgEnum_ChatMessageInfoTag)
	recvUserid := chatMessage.GetRecvUserid()
	roomId := chatMessage.GetRoomId()
	if recvUserid != 0 {
		recvPlayer := server.playerManager.findPlayerById(recvUserid)
		if recvPlayer != nil {
			recvPlayer.SendMessage(ChatMessageid, chatMessage)
		} else {
			res.Success = protobuf.ErrCodeType_PlayerOfflineError.Enum()
		}
	} else if roomId != 0 {
		room := server.roomManager.findRoom(uint64(roomId))
		if room != nil {
			room.sendMessageToPlayer(userid, ChatMessageid, chatMessage)
		} else {
			res.Success = protobuf.ErrCodeType_RoomNotExistError.Enum()
		}
	} else {
		res.Success = protobuf.ErrCodeType_ProtoPacketFormatError.Enum()
	}

	res.Timestamp = proto.Int32(int32(time.Now().Unix()))
	messageid := uint32(protobuf.MsgEnum_ChatResponseTag)
	player.SendMessage(messageid, res)
}
