package server

import (
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	go InitServer("./testdata/")

	var conn net.Conn
	for i := 0; i < 50; i++ { //Try connect to server. If you can't in 5s throw error
		var err error
		conn, err = net.Dial("tcp", "127.0.0.1:19151")
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	if conn == nil {
		t.Error("Cannot connect to server")
	}
}
