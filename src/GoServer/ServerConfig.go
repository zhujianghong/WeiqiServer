// ServerConfig
package main

import (
	"encoding/json"
	"os"
)

var Config GoConfig

func init() {
	Config.Load("../config/config.json")
	log.Println(Config)
}

type GoConfig struct {
	Addr                   string `json:"Addr"`
	PacketSendChanLimit    uint32 `json:"PacketSendChanLimit"`
	PacketReceiveChanLimit uint32 `json:"PacketReceiveChanLimit"`
	AcceptTimeout          uint32 `json:"AcceptTimeout"`
	RedisAddr              string `json:"RedisAddr"`
	RedisPasswd            string `json:"RedisPasswd"`
	MysqlName              string `json:"MysqlName"`
}

func (a *GoConfig) Load(path string) {
	fi, err := os.Stat(path)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	buff := make([]byte, fi.Size())
	f.Read(buff)
	f.Close()
	if err != nil {
		log.Fatal("path is error:", path, "error:", err)
	}
	err = json.Unmarshal(buff, &a)
	if err != nil {
		log.Fatal("json parse have a error:", err)
	}
}
