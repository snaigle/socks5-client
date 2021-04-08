package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	argLength := len(os.Args)
	if argLength < 3 {
		log.Fatal("cmd like: ./socks-client 127.0.0.1:8888 123.5.23.4")
		return
	}
	socksServer := os.Args[1]
	backend := os.Args[2]
	destAddr, err := ipstr2byte(backend)
	if err != nil {
		log.Fatal(err)
	}
	go listen(socksServer, destAddr, 443)
	go listen(socksServer, destAddr, 80)
	c := make(chan int)
	log.Println("server starting")
	<-c
}
func ipstr2byte(ip string) ([]byte, error) {
	args := strings.Split(ip, ".")
	if len(args) != 4 {
		return nil, errors.New("ip is not valid")
	}
	addr := make([]byte, 4)
	for i := 0; i < 4; i++ {
		v, e := strconv.Atoi(args[i])
		if e != nil {
			return nil, errors.New(fmt.Sprint(args[i], " can't parse to number"))
		}
		if v > 255 || v < 0 {
			return nil, errors.New(fmt.Sprint("ip byte ", v, " not valid"))
		}
		addr[i] = byte(v)
	}
	return addr, nil
}

func listen(socksServer string, destAddr []byte, port uint16) {
	address := ":" + strconv.Itoa(int(port))
	log.Println("start listen ", address, "to", destAddr)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		log.Println("listen ", address, "failed,reason:", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		go handleSocks5Connection(conn, socksServer, destAddr, port)
	}

}
func handleSocks5Connection(conn net.Conn, socksServer string, domain []byte, port uint16) {
	log.Println("accept conn at ", domain, ":", port)
	closed := false
	defer func() {
		if !closed {
			conn.Close()
		}
	}()
	socksConn, err := net.Dial("tcp", socksServer)
	if err != nil {
		log.Println("connect to socks server failed", err)
		return
	}
	// handshake
	_, err = socksConn.Write([]byte{0x05, 0x01, 0x00})
	if err != nil {
		log.Println("handshake failed", err)
		return
	}
	buf := make([]byte, 2)
	_, err = socksConn.Read(buf)
	portBuf := make([]byte, 2)
	binary.BigEndian.PutUint16(portBuf, port)
	send := []byte{0x05, 0x01, 0x00, 0x01}
	send = append(send, domain...)
	send = append(send, portBuf...)
	_, err = socksConn.Write(send)
	if err != nil {
		log.Println("write failed", err)
	}
	buf = make([]byte, 5)
	_, err = socksConn.Read(buf)
	if err != nil {
		log.Println("read failed", err)
	}
	buf = make([]byte, buf[4]+2)
	_, err = socksConn.Read(buf)
	if err != nil {
		log.Println("read failed", err)
	}
	go PipeThenClose(conn, socksConn)
	PipeThenClose(socksConn, conn)
	closed = true
	log.Println("closed connection to")
}
func PipeThenClose(src, dst net.Conn) {
	defer dst.Close()
	io.Copy(dst, src)
	return
}
