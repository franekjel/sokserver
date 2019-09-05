package server

import (
	"net"
	"os/exec"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	//first we copy test environment to /tmp to don't break it if something go wrong
	cmd := exec.Command("cp", "-r", "./testdata/", "/tmp/")
	cmd.Run()

	go InitServer("/tmp/testdata/")

	var conn net.Conn
	for i := 0; i < 50; i++ { //Try connect to server. If you can't in 5s throw error
		var err error
		conn, err = net.Dial("tcp", "127.0.0.1:19151") //It's easy to test sok on remote server - comment go InitServer() and set correct IP
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	if conn == nil {
		t.Error("Cannot connect to server")
	}
	time.Sleep(time.Millisecond * 300)
	t.Run("Create account", func(t *testing.T) {
		buff := []byte("login: Amandil\npassword: P@ssword\ncommand: create_account")
		n := len(buff)
		conn.SetDeadline(time.Now().Add(3 * time.Second))
		sendUint(conn, uint32(n))
		sendSlice(conn, buff)
		time.Sleep(time.Second)
		m := readUint(conn)
		buff = readNBytes(conn, m)
		if string(buff) != "status: ok\n" {
			t.Error("Bad return message: ", string(buff), "lenght:", len(buff), " - ", m)
		}
		conn.Close()
	})

}
