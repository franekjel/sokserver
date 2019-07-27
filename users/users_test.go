package users

import (
	"math/rand"
	"testing"
)

func BenchmarkGenHash(b *testing.B) {
	for i := 0; i < b.N; i++ {
		test := make([]byte, 32)
		rand.Read(test)
		genHash(test)
	}
}

func TestVerifyPassword(t *testing.T) {
	test := make([]byte, 32)
	rand.Read(test)
	back := make([]byte, 32)
	copy(back, test)
	salt, hash := genHash(test)
	u := User{PasswordHash: hash, PasswordSalt: salt}
	t.Log(test, back, u.PasswordSalt, u.PasswordHash)
	if !u.VerifyPassword(back) {
		t.Error(back, salt)
	}
}

func TestCheckGroup(t *testing.T) {
	grp := []string{"gr1", "gr2", "gr3"}
	u := User{Groups: grp}
	if !u.CheckGroup(&grp[1]) {
		t.Error(u.Groups, grp[1])
	}
	bad1 := "abcd"
	if u.CheckGroup(&bad1) {
		t.Error(u.Groups, bad1)
	}
	bad2 := "gr"
	if u.CheckGroup(&bad2) {
		t.Error(u.Groups, bad2)
	}
}
