// YFIMMain.go
package main

import (
	loger "YFTool/RpcLogs"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var log *loger.RpcLoger
var server *GoServer

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	startDate := time.Now().Unix()
	log = loger.NewRpcLogs("../log/GoServer-" + time.Now().Local().Format("2006-01-02") + "." + fmt.Sprintf("%d", startDate) + ".log")
	server = newGoServer()
	server.Start()
	log.Println("server start")
	chSig := make(chan os.Signal)
	signal.Notify(chSig, syscall.SIGINT, syscall.SIGTERM)
	log.Error("Signal: ", <-chSig)
	server.Stop()
}
