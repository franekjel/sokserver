package server

import (
	"net"
	"os"
	"os/exec"
	"testing"
	"time"
)

func connect() net.Conn {
	var conn net.Conn
	for i := 0; i < 30; i++ { //Try connect to server. If you can't in 3s throw error
		var err error
		conn, err = net.Dial("tcp", "127.0.0.1:19151")
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
	return conn
}

func createAccount(t *testing.T) {
	conn := connect()
	defer conn.Close()
	if conn == nil {
		t.Error("Cannot connect to server")
		return
	}
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

	//check if user file was created
	_, err := os.Stat("/tmp/testdata/users/Amandil.yml")
	if err != nil {
		t.Error("User file doesn't exist:", err.Error())
	}
}

func TestServer(t *testing.T) {
	//first we copy test environment to /tmp to don't break it if something go wrong
	cmd := exec.Command("rm", "-r", "/tmp/testdata/")
	cmd.Run()
	cmd = exec.Command("cp", "-r", "./testdata/", "/tmp/")
	cmd.Run()

	go InitServer("/tmp/testdata/")

	time.Sleep(time.Millisecond * 300)

	t.Run("Create account", createAccount)

}
