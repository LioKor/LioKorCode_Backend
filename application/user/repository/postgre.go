package repository

import (
	"liokoredu/application/user"
	"liokoredu/pkg/constants"
	"log"

	"github.com/gomodule/redigo/redis"
)

type UserDatabase struct {
	pool *redis.Pool
}

func (ud *UserDatabase) StoreSession(token string, uid uint64) error {
	client := ud.pool.Get()
	defer client.Close()

	_, err := client.Do("SET", token, uid, "EX", constants.WeekSec)
	if err != nil {
		log.Println("user repo: storeSession: error storing session")
		return err
	}

	return nil
}

func (ud *UserDatabase) GetId(token string) (uint64, error) {
	client := ud.pool.Get()
	defer client.Close()

	uid, err := client.Do("GET", token)
	if err != nil {
		log.Println("user repo: getId: error getting id")
		return 0, err
	}

	_, err = client.Do("EXPIRE", token, constants.WeekSec)
	if err != nil {
		log.Println("user repo: getId: error updating ttl")
		return uid.(uint64), err
	}

	return uid.(uint64), nil
}

func NewUserDatabase(pool *redis.Pool) user.Repository {
	return &UserDatabase{pool: pool}
}
