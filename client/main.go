package main

import (
	"net"
	"fmt"
	"github.com/fatih/color"
)

const server = "139.199.33.51:5000"
var q = make(chan(int))

func main() {
	conn, err := net.Dial("tcp", server)
	if err != nil {
		color.Red("connect error. cause => %s", err)
		return
	}
	color.Green("connect server [%s] success.", server)
	go read(conn)
	go waitForInput(conn)
	<- q
	color.Yellow("client close.")
}

func read(conn net.Conn) {
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			color.Red("read error. %s", err)
			q <- 1
			break
		}
		msg := string(data[0:n])
		color.Green("=> %s", msg)
	}
}

func waitForInput(conn net.Conn) {
	for {
		defer conn.Close()
		var msg string
		fmt.Scanln(&msg)
		if msg == ":q" {
			color.Red("quit...")
			q <- 1
			break
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			color.Red("client send error. %s", err)
			q <- 1
			break
		}
	}
}
