package models

import (
	"net/mail"
	"time"
)

type User struct {
	Id         uint64
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Fullname   string    `json:"fullname"`
	AvatarUrl  string    `json:"avatarUrl"`
	JoinedDate time.Time `json:"joinedDate"`
	Verified   bool      `json:"verified"`
	IsAdmin    bool      `json:"isAdmin"`
	JWT        string    `json:"jwtToken"`
}

type Avatar struct {
	AvatarUrl string `json:"avatarUrl"`
}

func (u *User) Validate() bool {
	if len(u.Username) < 6 || len(u.Username) > 30 {
		return false
	}

	if len(u.Email) == 0 || !isValidEmail(u.Email) {
		return false
	}

	if len(u.Password) < 6 || len(u.Password) > 30 {
		return false
	}

	if len(u.Fullname) < 6 || len(u.Fullname) > 50 {
		return false
	}

	return true
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

type UserAuth struct {
	Id       uint64
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserUpdate struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type PasswordNew struct {
	Old string `json:"oldPassword"`
	New string `json:"newPassword"`
}

func (u *UserUpdate) Validate() bool {
	if len(u.Email) == 0 || !isValidEmail(u.Email) {
		return false
	}

	if len(u.Fullname) < 6 || len(u.Fullname) > 50 {
		return false
	}

	return true
}

//easyjson:json
type Users []User
