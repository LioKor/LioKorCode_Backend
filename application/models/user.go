package models

import "time"

type User struct {
	Id         uint64
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	Fullname   string    `json:"email"`
	AvatarUrl  string    `json:"avatarUrl"`
	JoinedDate time.Time `json:"joinedDate"`
	IsAdmin    string    `json:"isAdmin"`
}

type UserAuth struct {
	Id       uint64
	Username string `json:"username"`
	Password string `json:"password"`
}

//easyjson:json
type Users []User
