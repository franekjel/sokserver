package users

import (
	"encoding/hex"
	"math/rand"

	"golang.org/x/crypto/argon2"
	"gopkg.in/yaml.v2"

	"github.com/franekjel/sokserver/fs"
	. "github.com/franekjel/sokserver/logger"
)

//User data
type User struct {
	Name         string   `yaml:"name"`
	Surname      string   `yaml:"surname"`
	PasswordHash string   `yaml:"password"`
	PasswordSalt string   `yaml:"salt"`
	Groups       []string `yaml:"groups,flow"`

	fs *fs.Fs
}

func (u *User) loadData() {
	if !u.fs.FileExist("user.yml") {
		Log(FATAL, "User data missing! %s", u.fs.Path)
	}
	buff := u.fs.ReadFile("user.yml")
	yaml.Unmarshal(buff, u)
}

func (u *User) saveData() {
	buff, err := yaml.Marshal(u)
	if err != nil {
		u.fs.WriteFile("user.yml", string(buff))
	}
}

//AddToGroup adds user to given group
func (u *User) AddToGroup(group *string) {
	u.Groups = append(u.Groups, *group)
}

func genHash(password []byte) (string, string) {
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		return "", ""
	}
	passwordHash := argon2.Key([]byte(password), salt, 3, 32*1024, 1, 32)
	return hex.EncodeToString(salt), hex.EncodeToString(passwordHash)
}

//VerifyPassword perform password verification based on stored hash and salt
func (u *User) VerifyPassword(password []byte) bool {
	salt, _ := hex.DecodeString(u.PasswordSalt)
	passwordHash := argon2.Key([]byte(password), salt, 3, 32*1024, 1, 32)
	if hex.EncodeToString(passwordHash) == u.PasswordHash {
		return true
	}
	return false
}

//CheckGroup checks if user is in given group
func (u *User) CheckGroup(group *string) bool {
	for _, s := range u.Groups {
		if s == *group {
			return true
		}
	}
	return false
}

//LoadUser create User based on data at given path
func LoadUser(path *fs.Fs) *User {
	user := new(User)
	user.fs = path
	user.loadData()
	Log(DEBUG, "name: %s, surname: %s, groups: %+q", user.Name, user.Surname, user.Groups)
	user.saveData()
	return user
}

//AddUser add new User. Function get path to "users" dir. Password will be hashed and erased
func AddUser(usersPath *fs.Fs, login *string, password []byte) *User {
	user := new(User)
	user.PasswordHash, user.PasswordSalt = genHash(password)
	for i := range password { //erase password for better security
		(password)[i] = 0
	}
	if user.PasswordHash == "" { //it means there is error in hashing
		return nil
	}

	usersPath.CreateDirectory(*login)
	user.fs = fs.Init(usersPath.Path, *login)
	user.saveData()
	Log(DEBUG, "Adding user: login: %s", *login)
	user.saveData()
	return user
}
