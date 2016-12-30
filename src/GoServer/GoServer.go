// GoServer
package main

import (
	"YFTool"
	"database/sql"
	"net"
	"protobuf"
	"sync"
	"time"

	"github.com/gansidui/gotcp"
	_ "github.com/go-sql-driver/mysql"
)

type GoServer struct {
	server        *gotcp.Server
	playerManager *PlayerManager
	roomManager   *RoomManager
	mysqlDB       *sql.DB
	matchService  *MatchService
	wg            sync.WaitGroup
	handlerFunMap map[protobuf.MsgEnum_MsgType]handlerFun
}

func newGoServer() *GoServer {
	s := &GoServer{
		playerManager: newPlayerManager(),
		roomManager:   newRoomManager(),
		matchService:  newMatchService(),
		handlerFunMap: make(map[protobuf.MsgEnum_MsgType]handlerFun),
	}

	NetConfig := &gotcp.Config{
		PacketSendChanLimit:    Config.PacketSendChanLimit,
		PacketReceiveChanLimit: Config.PacketReceiveChanLimit,
	}
	s.server = gotcp.NewServer(NetConfig, &GoCallBack{}, &GoProtocol{})
	var err error = nil
	s.mysqlDB, err = YFTool.GetMysqlDB(Config.MysqlName)
	if err != nil {
		log.Panic("err:", err)
	}

	return s
}

func (s *GoServer) Start() {
	s.registerHandler()
	tcpAddr, err := net.ResolveTCPAddr("tcp", Config.Addr)
	YFTool.CheckError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	YFTool.CheckError(err)
	go s.server.Start(listener, time.Duration(Config.AcceptTimeout))
	//开启匹配服务
	YFTool.AsyncDo(s.matchService.Start, &s.wg)
}

func (s *GoServer) Stop() {
	s.server.Stop()
	YFTool.CloseMysqlDB(s.mysqlDB)
	s.matchService.Stop()
	s.wg.Wait()
}

func (s *GoServer) registerHandler() {
	//服务器处理包的字典
	s.handlerFunMap[protobuf.MsgEnum_LoginRequestTag] = LoginHandler
	s.handlerFunMap[protobuf.MsgEnum_PlayerRegisterRequestTag] = PlayerRegisterHandler
	s.handlerFunMap[protobuf.MsgEnum_ChatRequestTag] = ChatHandler
	s.handlerFunMap[protobuf.MsgEnum_MatchGameRequestTag] = MatchGameHandler
	s.handlerFunMap[protobuf.MsgEnum_InviteFriendRequestTag] = InviteFriendHandler
	s.handlerFunMap[protobuf.MsgEnum_ReceiverFriendRequestTag] = ReceiverFriendHandler
	s.handlerFunMap[protobuf.MsgEnum_InvitePlayGameRequestTag] = InvitePlayGameHandler
	s.handlerFunMap[protobuf.MsgEnum_ReceiverPlayGameRequestTag] = ReceiverPlayGameHandler
	s.handlerFunMap[protobuf.MsgEnum_FriendListRequestTag] = FriendListHandler
	s.handlerFunMap[protobuf.MsgEnum_DeleteFriendRequestTag] = DeleteFriendHandler
}

func (s *GoServer) findHandler(cmd protobuf.MsgEnum_MsgType) handlerFun {
	if v, ok := s.handlerFunMap[cmd]; ok {
		return v
	}
	return nil
}
