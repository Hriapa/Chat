package model

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"time"
)

type User struct {
	Id                int
	Login             string `json:"login"`
	Password          string `json:"password,omitempty"`
	EncryptedPassword string `json:"-"`
}

type UserName struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Online bool   `json:"online"`
}

type UserInfo struct {
	Id         int    `json:"id"`
	Login      string `json:"login"`
	Name       string `json:"name"`
	Familyname string `json:"familyname"`
	Surname    string `json:"surname"`
	Birthdate  time.Time
}

type UserInList struct {
	Name   string
	Online bool
}

func (u *User) EncryptPassword() error {
	if len(u.Password) > 0 {
		pass := encryptString(u.Password)
		u.EncryptedPassword = pass
	}
	return nil
}

func encryptString(pass string) string {
	hash := md5.Sum([]byte(pass))
	return hex.EncodeToString(hash[:])
}

func (u *User) Private() {
	u.Password = ""
}

func (u *User) CheckPassword(pass string) bool {
	return encryptString(pass) == u.EncryptedPassword
}

func (u *User) CreateNewPassword() {
	digits := "0123456789"
	specials := "~=+%^*/()[]{}/!@#$?|"
	all := "ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		digits + specials
	length := 8
	buf := make([]byte, length)
	buf[0] = digits[rand.Intn(len(digits))]
	buf[1] = specials[rand.Intn(len(specials))]
	for i := 2; i < length; i++ {
		buf[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(buf), func(i, j int) {
		buf[i], buf[j] = buf[j], buf[i]
	})
	u.Password = string(buf)
}

func (u *UserInfo) CreateUserName() string {
	name := ""
	if u.Name != "" {
		name += u.Name
	}
	if u.Surname != "" {
		name += " " + u.Surname
	}
	if u.Familyname != "" {
		name += " " + u.Familyname
	}
	return name
}
