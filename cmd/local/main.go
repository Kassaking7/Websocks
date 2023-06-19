package main

import (
	"fmt"
	"github.com/Kassaking7/Websocks"
	"github.com/Kassaking7/Websocks/cmd"
	"log"
	"net"
)

const (
	DefaultListenAddr = ":7448"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)

	// default config
	config := &cmd.Config{
		ListenAddr: DefaultListenAddr,
	}
	config.ReadConfig()
	config.SaveConfig()

	// start local and listen
	lsLocal, err := Websocks.NewLsLocal(config.Password, config.ListenAddr, config.RemoteAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsLocal.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
Websocks-local:%s started，see the config:
Listening Address:
%s
Remote Address:
%s
Password:
%s`, version, listenAddr, config.RemoteAddr, config.Password))
	}))
}
