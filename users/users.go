package users

import (
	"encoding/hex"
	"math/rand"

	"golang.org/x/crypto/argon2"
	"gopkg.in/yaml.v2"

	"github.com/franekjel/sokserver/fs"
	log "github.com/franekjel/sokserver/logger"
)

//User data
type User struct {
	Login        string          `yaml:"-"`
	PasswordHash string          `yaml:"password"`
	PasswordSalt string          `yaml:"salt"`
	YamledGroups []string        `yaml:"groups"` //for keep in user file
	Groups       map[string]bool `yaml:"-"`
	fs           *fs.Fs
}

func (u *User) loadData() {
	if !u.fs.FileExist(u.Login + ".yml") {
		log.Fatal("User %s data missing! %s", u.Login, u.fs.Path)
	}
	buff := u.fs.ReadFile(u.Login + ".yml")
	yaml.Unmarshal(buff, u)
	u.Groups = make(map[string]bool, len(u.YamledGroups))
	for _, grp := range u.YamledGroups {
		u.Groups[grp] = true
	}
}

//SaveData save user config to file
func (u *User) SaveData() {
	buff, err := yaml.Marshal(u)
	if err == nil {
		u.fs.WriteFile(u.Login+".yml", string(buff))
	} else {
		log.Error("Cannot save user %s config", u.Login)
	}
}

//AddToGroup adds user to given group
func (u *User) AddToGroup(group string) {
	u.Groups[group] = true
	u.SaveData()
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
	if _, ok := u.Groups[*group]; ok {
		return true
	}
	return false
}

//LoadUser create User based on data at given path
func LoadUser(path *fs.Fs, login string) *User {
	user := User{Login: login, fs: path}
	user.loadData()
	log.Debug("name: %s, surname: %s, groups: %+q", user.Name, user.Surname, user.Groups)
	return &user
}

//AddUser add new User. Function get path to "users" dir. Password will be hashed and erased
func AddUser(usersPath *fs.Fs, login *string, password []byte) *User {
	user := new(User)
	user.fs = usersPath
	user.Login = *login
	user.PasswordHash, user.PasswordSalt = genHash(password)
	for i := range password { //erase password for better security
		(password)[i] = 0
	}
	if user.PasswordHash == "" { //it means there is error in hashing
		return nil
	}
	log.Debug("Adding user: login: %s", *login)
	user.SaveData()
	return user
}
