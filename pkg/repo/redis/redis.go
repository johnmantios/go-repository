package redis

import (
	"context"
	"database/sql"
	"errors"
	"github.com/redis/go-redis/v9"
	"johnmantios.com/go-repository/pkg/repo"
	"net/url"
	"os"
	"time"
)

type RedisRepo struct {
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

func NewRedisRepo(db *redis.Client) (*RedisRepo, error) {
	return &RedisRepo{
		Users: UsersModel{DB: db},
	}, nil
}

func (p RedisRepo) GetAUser(username string) (*repo.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	name, err := p.Users.DB.Get(ctx, username).Result()
	if err != nil {
		panic(err)
	}

	user := &repo.User{Name: name}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return &repo.User{}, ErrRecordNotFound
		default:
			return &repo.User{}, err
		}
	}

	return user, nil
}
