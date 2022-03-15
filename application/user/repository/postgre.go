package repository

import (
	"liokoredu/application/user"

	"github.com/gomodule/redigo/redis"
)

type UserDatabase struct {
	pool *redis.Pool
}

func NewUserDatabase(pool *redis.Pool) user.Repository {
	return &UserDatabase{pool: pool}
}
