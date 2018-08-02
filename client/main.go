package main

import (
	"net"
	"log"
	"fmt"
)

var q = make(chan(int))

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:5000")
	if err != nil {
		log.Fatalln("connect error.", err)
	}
	log.Println("connect success.")
	go read(conn)
	go waitForInput(conn)
	<- q
	log.Println("client close.")
}

func read(conn net.Conn) {
	defer conn.Close()
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			log.Println("read error.", err)
			q <- 1
			break
		}
		fmt.Println(string(data[0:n]))
	}
}

func waitForInput(conn net.Conn) {
	for {
		defer conn.Close()
		fmt.Printf("input(:q): \n")
		var msg string
		fmt.Scanln(&msg)
		if msg == ":q" {
			log.Println("quit...")
			q <- 1
			break
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Println("client send error.", err)
			q <- 1
			break
		}
	}
}
