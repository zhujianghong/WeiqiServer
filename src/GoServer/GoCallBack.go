// GoCallBack
package main

import (
	"protobuf"
	"sync"
	"sync/atomic"

	"github.com/gansidui/gotcp"
	"github.com/golang/protobuf/proto"
)

type handlerFun func(data []byte, conn *gotcp.Conn)

type GoCallBack struct {
}

type PlayerData struct {
	userid int32
	addr   string
	rwLock sync.RWMutex
}

func newPlayerData() *PlayerData {
	return &PlayerData{}
}

func (d *PlayerData) GetUserid() int32 {
	return atomic.LoadInt32(&d.userid)
}

func (this *GoCallBack) OnConnect(c *gotcp.Conn) bool {
	addr := c.GetRawConn().RemoteAddr().String()
	log.Println("OnConnect:", addr)
	c.PutExtraData(&PlayerData{addr: addr})
	return true
}

func (this *GoCallBack) OnMessage(c *gotcp.Conn, p gotcp.Packet) bool {
	packet := p.(*GoPacket)
	log.Println(packet)
	h := server.findHandler(protobuf.MsgEnum_MsgType(packet.cmd))
	if h != nil {
		h(packet.GetBody(), c)
	} else {
		log.Error("未处理的包[", packet.cmd, "]")
	}
	return true
}

func (this *GoCallBack) OnClose(c *gotcp.Conn) {
	log.Println("OnClose:", c.GetExtraData())
	playerData := c.GetExtraData().(*PlayerData)
	server.playerManager.logout(playerData.GetUserid(), c)
}

func sendMessageToConn(cmd uint32, message proto.Message, conn *gotcp.Conn) bool {
	data, err := proto.Marshal(message)
	if err != nil {
		log.Error(err)
		return false
	}
	packet := NewGoPacket(cmd, data)
	err = conn.AsyncWritePacket(packet, 0)
	if err != nil {
		log.Error(err)
		return false
	}
	return true
}
