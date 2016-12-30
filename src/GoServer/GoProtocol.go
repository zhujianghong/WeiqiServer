// ServerLogic
package main

import (
	"encoding/binary"
	"io"
	"net"
	"time"

	"github.com/gansidui/gotcp"
)

type GoProtocol struct {
}

func (p *GoProtocol) ReadPacket(conn *net.TCPConn) (gotcp.Packet, error) {
	var (
		lengthBytes []byte = make([]byte, 4)
		length      uint32
	)
	timer := time.AfterFunc(60*time.Second, func() {
		conn.Close()
	})
	if _, err := io.ReadFull(conn, lengthBytes); err != nil {
		log.Error("io.ReadFull ", err)
		return nil, err
	}
	log.Println("headerData:", lengthBytes)
	length = binary.BigEndian.Uint32(lengthBytes)

	buff := make([]byte, length)
	copy(buff[0:4], lengthBytes)
	_, err := io.ReadFull(conn, buff[4:])
	log.Println("Data:", buff)
	if timer != nil {
		timer.Stop()
	}
	if err != nil {
		log.Error(err)
		return nil, err
	}
	packet := NewGoPacket(0, make([]byte, 0))
	err = packet.Parse(buff)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return packet, nil
}
