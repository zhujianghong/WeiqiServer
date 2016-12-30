// ServerLogic
package main

import (
	"Aes"
	//	"Tea"
	"encoding/binary"
)

type GoPacket struct {
	length uint32
	cmd    uint32
	buff   []byte
}

func NewGoPacket(cmd uint32, data []byte) *GoPacket {
	p := &GoPacket{}
	if len(data) != 0 && cmd != 0 {
		p.length = uint32(8 + len(data))
		p.cmd = cmd
		p.buff = data
	}
	return p
}

func (p *GoPacket) GetLength() uint32 {
	return p.length
}

func (p *GoPacket) GetBody() []byte {
	return p.buff
}

func (p *GoPacket) GetCmd() uint32 {
	return p.cmd
}

func (p *GoPacket) Serialize() []byte {

	enBData := make([]byte, p.length-4)
	binary.BigEndian.PutUint32(enBData[0:4], p.cmd)
	copy(enBData[4:], p.buff)
	enAData, _ := Aes.Encrypt(enBData, Aes.AESKEY)
	length := len(enAData) + 4
	data := make([]byte, length)
	binary.BigEndian.PutUint32(data[0:4], uint32(length))
	copy(data[4:], enAData)
	log.Println(data)

	/*dataLen := Aes.GetEncryptAfterLen(int32(p.length-4)) + 4
	data := make([]byte, dataLen)
	binary.BigEndian.PutUint32(data[0:4], uint32(dataLen))
	binary.BigEndian.PutUint32(data[4:8], p.cmd)
	copy(data[8:], p.buff)
	enData, _ := Aes.Encrypt(data[4:len(p.buff)+4], Aes.AESKEY)
	copy(data[4:], enData)
	log.Println(data)*/

	/*data := make([]byte, p.GetLength())
	binary.BigEndian.PutUint32(data[0:4], p.length)
	binary.BigEndian.PutUint32(data[4:8], p.cmd)
	copy(data[8:], p.buff)
	copy(data[4:], Tea.Encrypt(data, 4, len(data), 16))*/
	return data
}

func (p *GoPacket) Parse(data []byte) error {
	deData, err := Aes.Decrypt(data[4:], Aes.AESKEY)
	log.Println(deData)
	if err != nil {
		return err
	}
	p.length = uint32(len(deData) + 4)
	p.buff = make([]byte, p.length-8)
	p.cmd = binary.BigEndian.Uint32(deData[0:4])
	copy(p.buff, deData[4:])
	return nil
	/*p.length = uint32(len(data))
	p.buff = make([]byte, p.length-8)
	deData := Tea.Decrypt(data[4:], 0, len(data)-4, 16)
	p.cmd = binary.BigEndian.Uint32(deData[0:4])
	copy(p.buff, deData[4:])*/
}
