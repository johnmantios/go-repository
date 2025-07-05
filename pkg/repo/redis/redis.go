package redis

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"johnmantios.com/go-repository/pkg/service"
	"net/url"
	"os"
	"time"
)

type Repo struct {
	Users UsersModel
}

type UsersModel struct {
	DB *redis.Client
}

func OpenDB() (*redis.Client, error) {
	present := false
	password, present := os.LookupEnv("DB_PASSWORD")
	if !present {
		return nil, errors.New("env variable DB_PASSWORD missing")
	}
	password = url.QueryEscape(password) //escaping in case of weird password

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: password, // no password set
		DB:       0,        // use default DB
	})

	return rdb, nil
}

func NewRedisRepo(db *redis.Client) (*Repo, error) {
	return &Repo{
		Users: UsersModel{DB: db},
	}, nil
}

func (p Repo) GetAUser(username string) (*service.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name, err := p.Users.DB.Get(ctx, username).Result()
	if err != nil {
		panic(err)
	}

	user := &service.User{Name: name}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &service.User{}, ErrRecordNotFound
		default:
			return &service.User{}, err
		}
	}

	return user, nil
}
