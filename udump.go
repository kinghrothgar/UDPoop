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

	if err := conf.Parse(); err != nil {
		gslog.Fatal("MAIN: failed to parse conf with error: %s", err.Error())
	}
	addr, err := net.ResolveUDPAddr("udp4", "127.0.0.1:60000")
	if err != nil {
		gslog.Fatal(err.Error())
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		gslog.Fatal(err.Error())
	}
	defer conn.Close()
	r := bufio.NewReaderSize(conn, 52428800)
	f, err := os.OpenFile("/tmp/udpoop", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0664)
	if err != nil {
		gslog.Fatal(err.Error())
	}
	defer f.Close()
	w := bufio.NewWriterSize(f, 52428800)
	go func() {
		n, err := w.ReadFrom(r)
		if err != nil {
			gslog.Fatal(err.Error())
		}
		gslog.Info("Read %d bytes", n)
	}()
	tiktok := time.NewTicker(5 * time.Second)
	for {
		<-tiktok.C
		w.Flush()
	}
}
