package main

import (
	"bufio"
	"github.com/grooveshark/golib/gslog"
	"github.com/kinghrothgar/UDumP/conf"
	"net"
	"os"
	"time"
)

// To be set at build
var (
	buildCommit string
	buildDate   string
)

func main() {
	gslog.Info("Goblin started [build commit: %s, build date: %s]", buildCommit, buildDate)

	// Setup config and logging
	if err := conf.Parse(); err != nil {
		gslog.Fatal("MAIN: failed to parse conf with error: %s", err.Error())
	}
	gslog.SetMinimumLevel(conf.GetStr("--log-level"))
	if logFile, set := conf.ParamStr("--log-file"); set {
		gslog.SetLogFile(logFile)
	}

	addr, err := net.ResolveUDPAddr("udp4", conf.GetStr("--host")+":"+conf.GetStr("--port"))
	if err != nil {
		gslog.Fatal(err.Error())
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		gslog.Fatal(err.Error())
	}
	defer conn.Close()
	r := bufio.NewReaderSize(conn, conf.GetInt("--buffer"))
	f, err := os.OpenFile("/tmp/udpoop", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		gslog.Fatal(err.Error())
	}
	defer f.Close()
	w := bufio.NewWriterSize(f, conf.GetInt("--buffer"))
	go func() {
		n, err := w.ReadFrom(r)
		if err != nil {
			gslog.Fatal(err.Error())
		}
		gslog.Info("Read %d bytes", n)
	}()
	tiktok := time.NewTicker(time.Duration(conf.GetInt("--flush")) * time.Second)
	for {
		<-tiktok.C
		w.Flush()
	}
}
