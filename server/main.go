package main

import (
	"fmt"
	"bytes"
	"log"
	"net"
	"sync"
)

// 连接管理
type ConnM struct {
	sync.RWMutex
	m map[string]net.Conn
}

// 添加连接
func (cm *ConnM) Add(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	cm.Lock()
	defer cm.Unlock()
	cm.m[addr] = conn
}

// 移除连接
func (cm *ConnM) Remove(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	cm.Lock()
	defer cm.Unlock()
	delete(cm.m, addr)
}

func (cm *ConnM) Boradcast(conn net.Conn, msg []byte) {
	cm.RLock()
	defer cm.RUnlock()
	for _, c := range cm.m {
		if c == conn {
			continue
		}
		_, err := c.Write(msg)
		if err != nil {
			log.Println("write error", err)
		}
	}
}

func (cm *ConnM) Len() int {
	cm.RLock()
	defer cm.RUnlock()
	return len(cm.m)
}

func NewConnM() *ConnM {
	return &ConnM{
		m: make(map[string]net.Conn),
	}
}

var cm = NewConnM()

func main() {
	ln, err := net.Listen("tcp", ":5000")
	if err != nil {
		log.Fatalln(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("server accept error.", err)
		}
		join(conn)
		go handleConnection(conn)
	}
}

// 处理连接
func handleConnection(conn net.Conn) {
	defer func() {
		leave(conn)
		conn.Close()
	}()
	// 数据缓存
	data := make([]byte, 1024)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println("conn read error.", err)
			break
		}
		b := data[0:n]

		addr := conn.RemoteAddr().String()

		var buf bytes.Buffer
		buf.WriteString(addr)
		buf.WriteString(": ")
		buf.Write(b)
		cm.Boradcast(conn, buf.Bytes())
	}
}

// 加入
func join(conn net.Conn) {
	addr := conn.RemoteAddr().String()
	log.Printf("[%s] join\n", addr)
	msg := fmt.Sprintf("welcome, [%d] person online.", cm.Len() + 1)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Println(err)
		return
	}
	cm.Boradcast(conn, []byte("sys msg [" + addr + "] join"))
	cm.Add(conn)
}

// 离开
func leave(conn net.Conn) {
	cm.Remove(conn)
	addr := conn.RemoteAddr().String()
	log.Printf("[%s] leave\n", addr)
	cm.Boradcast(conn, []byte("sys msg [" + addr + "] leave"))
}
