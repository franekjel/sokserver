package server

import (
	"encoding/binary"
	"net"
	"strconv"
	"time"

	log "github.com/franekjel/sokserver/logger"
)

type connectionData struct {
	conn net.Conn
	data []byte
}

//reads 4 bytes and convert it to uint
func readUint(conn net.Conn) uint {
	buff := readNBytes(conn, 4)
	return uint(binary.BigEndian.Uint32(buff))
}

//reads n bytes from conn
func readNBytes(conn net.Conn, n uint) []byte {
	buff := make([]byte, n)
	var i uint
	for i < n {
		nbyte, err := conn.Read(buff[i:n])
		i += uint(nbyte)
		if err != nil {
			log.Warn("Cannot read message: %s", err.Error())
			return nil
		}
	}
	return buff
}

//send given uint
func sendUint(conn net.Conn, n uint32) {
	buff := make([]byte, 4)
	binary.BigEndian.PutUint32(buff, n)
	var i int
	for i < 4 {
		n, err := conn.Write(buff[i:4])
		i += n
		if err != nil {
			log.Error("Cannot send data to client %s. %s", conn.RemoteAddr().String(), err.Error())
			return
		}
	}
}

//send given slice
func sendSlice(conn net.Conn, buff []byte) {
	var i = 0
	n := len(buff)
	for i < n {
		nbyte, err := conn.Write(buff[i:n])
		i += nbyte
		if err != nil {
			log.Error("Cannot send slice to client")
			return
		}
	}
}

//reads data from client and send to chan
func readMessage(conn net.Conn, ch chan *connectionData) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	msgSize := readUint(conn)
	if msgSize == 0 || msgSize > 128*1024 {
		log.Error("Bad message size or other error from %s", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	buff := readNBytes(conn, msgSize)
	if buff == nil {
		log.Error("Cannot get data from %s", conn.RemoteAddr().String())
		conn.Close()
		return
	}
	data := connectionData{
		conn,
		buff,
	}
	ch <- &data
}

func (s *Server) startListening(ch chan *connectionData) {
	l, err := net.Listen("tcp", ":"+strconv.FormatUint(uint64(s.conf.Port), 10))
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

func sendResponse(conn net.Conn, buff []byte) {
	sendUint(conn, uint32(len(buff)))
	sendSlice(conn, buff)
	conn.Close()
}
