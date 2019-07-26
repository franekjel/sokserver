package users

import(
	. "github.com/franekjel/sokserver/logger"
)

type User struct {
	login string
	name string
	surname string
	passwordHash string
	groups []string
}

func NewUser(dir string) *User{
	user := new(User)
	Log(DEBUG, "login: %s, name: %s, surname %s", user.login, user.name, user.surname)
	return user
}
