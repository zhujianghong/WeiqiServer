// LoginLogicHandler
package main

import (
	"protobuf"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

func LoginHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.LoginRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)

	res := &protobuf.LoginResponse{}
	resMessageid := protobuf.MsgEnum_LoginResponseTag
	res.Success = protobuf.ErrCodeType_Success.Enum()
	info, errCode := PlayerLogin(req.GetAccount(), req.GetPasswd())
	if errCode == 0 {
		player := newPlayer(info.Userid, req.GetAccount(), conn)
		player.name = info.Name
		player.sex = info.Sex
		player.grade = info.Grade
		player.icon = info.Icon
		player.money = info.Money
		player.autograph = info.Autograph
		player.gnum = info.Gnum
		player.gsnum = info.Gsnum
		player.gfnum = info.Gfnum
		player.gdnum = info.Gdnum
		player.fgsnum = info.Fgsnum
		player.fgfnum = info.Fgfnum
		player.fgdnum = info.Fgdnum
		flag := server.playerManager.login(player)
		if flag {
			res.Info.Userid = proto.Int32(info.Userid)
			res.Info.Name = proto.String(info.Name)
			res.Info.Icon = proto.Int32(info.Icon)
			res.Info.Sex = proto.Int32(info.Sex)
			res.Info.Autograph = proto.String(info.Autograph)
			res.Info.Grade = proto.Int32(info.Grade)
			res.Info.Gnum = proto.Int32(info.Gnum)
			res.Info.Gsnum = proto.Int32(info.Gsnum)
			res.Info.Gfnum = proto.Int32(info.Gfnum)
			res.Info.Gdnum = proto.Int32(info.Gdnum)
			res.Info.Fgsnum = proto.Int32(info.Fgsnum)
			res.Info.Fgfnum = proto.Int32(info.Fgfnum)
			res.Info.Fgdnum = proto.Int32(info.Fgdnum)
			res.Info.Money = proto.Int32(info.Money)
			player.SendMessage(uint32(resMessageid), res)
			return
		} else {
			res.Success = protobuf.ErrCodeType_RepeatedLoginError.Enum()
		}
	} else {
		res.Success = errCode.Enum()
	}
	sendMessageToConn(uint32(resMessageid), res, conn)
}

func PlayerRegisterHandler(data []byte, conn *gotcp.Conn) {
	log.Println(data)
	req := &protobuf.PlayerRegisterRequest{}
	err := proto.Unmarshal(data, req)
	if err != nil {
		log.Error("parse packet err", err)
		return
	}
	log.Println(req)
	userid, errCode := AddRegisterPlayerInfo(req)
	res := &protobuf.PlayerRegisterResponse{}
	res.Success = errCode.Enum()
	res.Userid = proto.Int32(userid)
	log.Println(res.GetSuccess(), res.GetUserid())
	resMessageid := uint32(protobuf.MsgEnum_PlayerRegisterResponseTag)
	sendMessageToConn(resMessageid, res, conn)
}
