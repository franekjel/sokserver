package server

import (
	"bufio"
	"encoding/binary"
	"net"
	"strconv"
	"time"

	log "github.com/franekjel/sokserver/logger"
)

//reads 4 bytes and convert it to uint
func readUint(conn *net.Conn) uint {
	buff := make([]byte, 4)
	c := bufio.NewReader(*conn)
	for i := 0; i < 4; i++ {
		b, err := c.ReadByte()
		if err != nil {
			return 0
		}
		buff[i] = b
	}
	return uint(binary.BigEndian.Uint32(buff))
}

//reads n bytes from conn
func readNBytes(conn *net.Conn, n uint) []byte {
	buff := make([]byte, n)
	var i uint
	for {
		nbyte, err := (*conn).Read(buff[i:])
		if err != nil {
			return nil
		}
		i += uint(nbyte)

		if i == n {
			return buff
		}
	}
}

//reads data from client and send to chan
func readMessage(conn net.Conn, ch chan *[]byte) {
	defer conn.Close()
	conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	msgSize := readUint(&conn)
	if msgSize == 0 || msgSize > 128*1024 {
		log.Error("Bad message size or other error from %s", conn.RemoteAddr().String())
		return
	}
	buff := readNBytes(&conn, msgSize)
	if buff == nil {
		log.Error("Cannot get data from %s", conn.RemoteAddr().String())
		return
	}
	ch <- &buff
}

func (s *Server) startListening(ch chan *[]byte) {
	l, err := net.Listen("tcp", "127.0.0.1:"+strconv.FormatUint(uint64(s.conf.Port), 10))
	if err != nil {
		log.Fatal("Cannot start listening on port %s: %s", strconv.FormatUint(uint64(s.conf.Port), 10), err.Error())
	}
	defer l.Close()
	log.Info("Starting listening on port %d", s.conf.Port)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error("Cannot accept connection: %s", err.Error())
			continue
		}
		go readMessage(conn, ch)
	}
}
