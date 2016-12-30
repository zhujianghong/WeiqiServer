// Player
package main

import (
	"sync/atomic"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

var (
	idle    int32 = 0
	playing int32 = 1
	finish  int32 = 2
)

type Player struct {
	userid    int32  //玩家唯一id
	account   string //玩家的账号
	name      string //玩家昵称
	grade     int32  //玩家段位
	money     int32  //玩家钱币
	icon      int32  //玩家头像id
	sex       int32  //玩家性别
	autograph string //玩家签名
	roomId    uint64 //玩家所在的房间号
	status    int32  //是否在游戏中的状态
	gnum      int32  //总局数
	gsnum     int32  //升降局胜利的局数
	gfnum     int32  //升降局失败的局数
	gdnum     int32  //升降局平局的局数
	fgsnum    int32  //友谊局胜利的局数
	fgfnum    int32  //友谊局失败的局数
	fgdnum    int32  //友谊局平局的局数
	conn      *gotcp.Conn
}

func newPlayer(userid int32, account string, conn *gotcp.Conn) *Player {
	return &Player{
		userid:  userid,
		account: account,
		status:  idle,
		conn:    conn,
	}
}

func (p *Player) SendMessage(cmd uint32, message proto.Message) bool {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Error(err)
		return false
	}
	packet := NewGoPacket(cmd, data)
	err = p.conn.AsyncWritePacket(packet, 0)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}

func (p *Player) play() {
	atomic.StoreInt32(&p.status, playing)
}

func (p *Player) finishGame() {
	atomic.StoreInt32(&p.status, finish)
}

func (p *Player) quitRoom() {
	atomic.StoreInt32(&p.status, idle)
}

func (p *Player) isPlaying() bool {
	return atomic.LoadInt32(&p.status) == playing
}

func (p *Player) isIdle() bool {
	return atomic.LoadInt32(&p.status) == idle
}
