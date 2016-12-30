package main

import (
	"YFTool"
	log "YFTool/RpcLogs"
	"crypto/md5"
	"encoding/binary"
	"net"
	"protobuf"
	"tea"

	"github.com/golang/protobuf/proto"
	//"github.com/gansidui/gotcp"
)

func main() {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "192.168.2.171:8000")
	YFTool.CheckError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	YFTool.CheckError(err)
	log.Println("success")

	req := &protobuf.LoginRequest{}
	req.Account = proto.String("zhu")
	passwd := "123456"
	req.Kind = proto.Int(1)
	md5Handler := md5.New()
	md5Handler.Write([]byte(passwd))
	data := md5Handler.Sum(nil)
	req.Passwd = proto.String(string(data))

	reqData, err := proto.Marshal(req)
	if err != nil {
		log.Println(err)
	}

	length := uint32(len(reqData) + 8)
	sendData := make([]byte, length)
	binary.BigEndian.PutUint32(sendData[0:4], length)
	cmd := uint32(protobuf.MsgEnum_LoginRequestTag)
	binary.BigEndian.PutUint32(sendData[4:8], cmd)
	copy(sendData[8:], reqData)
	copy(sendData[4:], Tea.Encrypt(sendData, 4, len(sendData), 16))

	if _, err := conn.Write(sendData); err != nil {
		log.Error("err:", err)
		return
	}
	log.Println("Write success")

	for {

	}

	conn.Close()
}
