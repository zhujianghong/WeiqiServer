// FriendLogicHandler
package main

import (
	"protobuf"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

func InviteFriendHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.InviteFriendRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	userid := conn.GetExtraData().(*PlayerData).GetUserid()
	player := server.playerManager.findPlayerById(userid)
	if player == nil {
		conn.Close()
		return
	}
	log.Println(req)
	res := &protobuf.InviteFriendResponse{}
	resMesageid := uint32(protobuf.MsgEnum_InviteFriendResponseTag)
	res.Success = protobuf.ErrCodeType_Success.Enum()
	inviteFriendInfo := req.Info
	messageid := uint32(protobuf.MsgEnum_InviteFriendInfoTag)
	recvUserid := inviteFriendInfo.GetReceiverUserid()
	recvPlayer := server.playerManager.findPlayerById(recvUserid)
	if recvPlayer != nil {
		code := InviteOrAccpetFriend(inviteFriendInfo.GetInviterUserid(), inviteFriendInfo.GetReceiverUserid(), 0)
		if code == -1 {
			res.Success = protobuf.ErrCodeType_MysqlError.Enum()
		} else if code == 0 {
			recvPlayer.SendMessage(messageid, inviteFriendInfo)
		} else if code == 1 {
			res.Success = protobuf.ErrCodeType_InviterFriendListLimitError.Enum()
		} else if code == 2 {
			res.Success = protobuf.ErrCodeType_InviteeFriendListLimitError.Enum()
		}
	} else {
		res.Success = protobuf.ErrCodeType_PlayerOfflineError.Enum()
	}
	player.SendMessage(resMesageid, res)
}

func ReceiverFriendHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.ReceiverFriendRequest{}
	res := &protobuf.ReceiverFriendResponse{}
	resMessageid := uint32(protobuf.MsgEnum_ReceiverFriendResponseTag)
	res.Success = protobuf.ErrCodeType_Success.Enum()
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	player := server.playerManager.findPlayerById(req.GetUserid())
	if player == nil {
		conn.Close()
		return
	}

	recvInfo := req.GetInfo()
	if recvInfo.GetIsAgree() {
		// 判断此玩家是否加满好友了 还要判断对方是否好友已满
		code := InviteOrAccpetFriend(recvInfo.GetInviterUserid(), recvInfo.GetReceiverUserid(), 1)
		if code == -1 {
			res.Success = protobuf.ErrCodeType_MysqlError.Enum()
		} else if code == 0 {
			inviterPlayer := server.playerManager.findPlayerById(req.GetInfo().GetInviterUserid())
			inviterMessageid := uint32(protobuf.MsgEnum_ReceiverFriendInfoTag)
			if inviterPlayer != nil {
				inviterPlayer.SendMessage(inviterMessageid, recvInfo)
			}
		} else if code == 1 {
			res.Success = protobuf.ErrCodeType_InviterFriendListLimitError.Enum()
		} else if code == 2 {
			res.Success = protobuf.ErrCodeType_InviteeFriendListLimitError.Enum()
		}
	}
	player.SendMessage(resMessageid, res)
}

func DeleteFriendHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.DeleteFriendRequest{}
	res := &protobuf.DeleteFriendResponse{}
	res.Success = protobuf.ErrCodeType_Success.Enum()
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	player := server.playerManager.findPlayerById(req.GetUserid())
	if player == nil {
		conn.Close()
		return
	}
	if DeleteFriend(req.GetUserid(), req.GetInfo().GetFriendUserid()) {
		friendPlayer := server.playerManager.findPlayerById(req.GetInfo().GetFriendUserid())
		if friendPlayer != nil {
			Deletemessageid := uint32(protobuf.MsgEnum_DeleteFriendInfoTag)
			friendPlayer.SendMessage(Deletemessageid, req.GetInfo())
		}
	} else {
		res.Success = protobuf.ErrCodeType_MysqlError.Enum()
	}

	messageid := uint32(protobuf.MsgEnum_DeleteFriendResponseTag)
	player.SendMessage(messageid, res)
}

func FriendListHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.FriendListRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	player := server.playerManager.findPlayerById(req.GetUserid())
	if player == nil {
		conn.Close()
		return
	}

	res := GetFriendList(player.userid)
	messageid := uint32(protobuf.MsgEnum_FriendListResponseTag)
	player.SendMessage(messageid, res)
}
