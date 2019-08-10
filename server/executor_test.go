package server

import (
	"testing"

	"github.com/franekjel/sokserver/fs"
	"github.com/franekjel/sokserver/users"
)

func TestCreateAccountOK(t *testing.T) {
	//simple call function in /tmp, should work
	var s Server
	s.fs = fs.Init("/tmp", "")
	s.fs.CreateDirectory("users")
	s.users = make(map[string]*users.User)
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "create_account",
	}
	buff := s.createAccount(&com)
	if string(buff) != "status: ok\n" {
		t.Error(string(buff))
	}
}

func TestCreateAccountDuplicated(t *testing.T) {
	//create account two times
	var s Server
	s.fs = fs.Init("/tmp", "")
	s.fs.CreateDirectory("users")
	s.users = make(map[string]*users.User)
	com := Command{
		Login:    "testLogin",
		Password: "password",
		Command:  "create_account",
	}
	buff := s.createAccount(&com)
	buff = s.createAccount(&com)
	if string(buff) != "status: Cannot create account, there is already user with this login\n" {
		t.Error(string(buff))
	}
}
