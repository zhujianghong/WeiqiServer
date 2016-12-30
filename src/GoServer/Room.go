// Room
package main

import (
	"protobuf"
	"sync"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

type Room struct {
	roomId      uint64
	blackUserid int32
	whiteUserid int32
	PlayerMap   map[int32]*Player
	rwLock      sync.RWMutex
}

func newRoom(roomId uint64) *Room {
	r := &Room{
		roomId:      roomId,
		blackUserid: 0,
		whiteUserid: 0,
		PlayerMap:   make(map[int32]*Player),
	}
	return r
}

func (r *Room) initRoom(send *Player, recv *Player) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	r.PlayerMap[send.userid] = send
	r.PlayerMap[recv.userid] = recv

}

func (r *Room) beginGame() {
	r.setPlayerColor()
}

func (r *Room) destroyRoom() {
	//通知房间的玩家房间销毁
	server.roomManager.delRoom(r.roomId)
}

func (r *Room) broadcastMessage(userid int32, cmd uint32, message proto.Message) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	for i, v := range r.PlayerMap {
		if i != userid {
			v.SendMessage(cmd, message)
		}
	}
}

func (r *Room) sendMessageToPlayer(userid int32, cmd uint32, message proto.Message) {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if v, ok := r.PlayerMap[userid]; ok {
		v.SendMessage(cmd, message)
	}
}

func (r *Room) FindPlayerByID(userid int32) *Player {
	r.rwLock.RLock()
	defer r.rwLock.RUnlock()
	if v, ok := r.PlayerMap[userid]; ok {
		return v
	}
	return nil
}

func (r *Room) quitRoom(userid int32, conn *gotcp.Conn) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	if v, ok := r.PlayerMap[userid]; ok {
		if v.conn == conn {
			delete(r.PlayerMap, userid)
			//TODO 通知其他玩家 userid离场
		}
	}
}

func (r *Room) enterRoom(p *Player) {
	r.rwLock.Lock()
	defer r.rwLock.Unlock()
	if v, ok := r.PlayerMap[p.userid]; ok {
		if v.conn != p.conn {
			delete(r.PlayerMap, p.userid)
			//TODO 通知其他玩家 userid离场
		} else {
			return
		}
	}
	r.PlayerMap[p.userid] = p
	//TODO 通知其他玩家 userid进入房间
}

func (r *Room) SendMatchSuccessInfo() {
	info := &protobuf.MatchGameInfo{}
	info.Success = protobuf.ErrCodeType_Success.Enum()
	messageid := uint32(protobuf.MsgEnum_MatchGameInfoTag)
	info.WhiteUserid = proto.Int32(r.whiteUserid)
	info.BlackUserid = proto.Int32(r.blackUserid)

	whitePlayer := r.FindPlayerByID(r.whiteUserid)
	blackPlayer := r.FindPlayerByID(r.blackUserid)
	if whitePlayer != nil && blackPlayer != nil {
		info.OtherInfo = CreatePlayerInfo(blackPlayer)
		whitePlayer.SendMessage(messageid, info)
		info.OtherInfo = CreatePlayerInfo(whitePlayer)
		blackPlayer.SendMessage(messageid, info)
	} else {
		info.Success = protobuf.ErrCodeType_PlayerLeaveRoomError.Enum()
		r.sendMessageToPlayer(0, messageid, info)
		r.destroyRoom()
	}
}

func (r *Room) setPlayerColor() {
	r.rwLock.RLock()
	for k, _ := range r.PlayerMap {
		if r.whiteUserid == 0 {
			r.whiteUserid = k
		} else {
			r.blackUserid = k
		}
	}
	r.rwLock.RUnlock()
}

func CreatePlayerInfo(player *Player) *protobuf.PlayerInfo {
	playerInfo := &protobuf.PlayerInfo{}
	playerInfo.Userid = proto.Int32(player.userid)
	playerInfo.Autograph = proto.String(player.autograph)
	playerInfo.Fgdnum = proto.Int32(player.fgdnum)
	playerInfo.Fgfnum = proto.Int32(player.fgfnum)
	playerInfo.Fgsnum = proto.Int32(player.fgsnum)
	playerInfo.Gdnum = proto.Int32(player.gdnum)
	playerInfo.Gfnum = proto.Int32(player.gfnum)
	playerInfo.Gnum = proto.Int32(player.gnum)
	playerInfo.Grade = proto.Int32(player.grade)
	playerInfo.Gsnum = proto.Int32(player.gsnum)
	playerInfo.Icon = proto.Int32(player.icon)
	playerInfo.Money = proto.Int32(player.money)
	playerInfo.Name = proto.String(player.name)
	playerInfo.Sex = proto.Int32(player.sex)
	return playerInfo
}
