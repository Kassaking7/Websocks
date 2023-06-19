package main

import (
	"fmt"
	"github.com/Kassaking7/Websocks"
	"github.com/Kassaking7/Websocks/cmd"
	"github.com/phayes/freeport"
	"log"
	"net"
	"os"
	"strconv"
)

var version = "master"

func main() {
	log.SetFlags(log.Lshortfile)
	port, err := strconv.Atoi(os.Getenv("WEBSOCKS_SERVER_PORT"))
	if err != nil {
		port, err = freeport.GetFreePort()
	}
	if err != nil {
		port = 7448
	}
	// default config
	config := &cmd.Config{
		ListenAddr: fmt.Sprintf(":%d", port),
		Password:   Websocks.RandPassword(),
	}
	config.ReadConfig()
	config.SaveConfig()

	// start server and listen
	lsServer, err := Websocks.NewLsServer(config.Password, config.ListenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	log.Fatalln(lsServer.Listen(func(listenAddr *net.TCPAddr) {
		log.Println(fmt.Sprintf(`
Websocks-server:%s started，see the config:
Listening Address:
%s
Password:
%s`, version, listenAddr, config.Password))
	}))
}
