package redis

import (
	"context"
	"database/sql"
	"errors"
	"johnmantios.com/go-repository/pkg/repo"
	"time"
)

type RedisRepo struct {
	Users UsersModel
}

type UsersModel struct {
	DB *sql.DB
}

func NewRedisRepo(db *sql.DB) (*RedisRepo, error) {
	return &RedisRepo{
		Users: UsersModel{DB: db},
	}, nil
}

func (p RedisRepo) GetAUser(username string) (repo.User, error) {
	query := `
		SELECT 
		    name
		FROM 
		    repo.user
		WHERE name = $1
		`

	var user repo.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rowContext := p.Users.DB.QueryRowContext(ctx, query, username)

	err := rowContext.Scan(
		&user.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return repo.User{}, ErrRecordNotFound
		default:
			return repo.User{}, err
		}
	}

	return user, nil
}
