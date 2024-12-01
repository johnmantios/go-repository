package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"johnmantios.com/go-repository/pkg/repo"
	"net/url"
	"os"
	"time"
)

type PostgresRepo struct {
	Users UsersModel
}

type UsersModel struct {
	DB *sql.DB
}

type Db struct {
	Username     string
	Password     string
	Host         string
	Name         string
	Ssl          string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  string
}

type Config struct {
	Env  string
	Port int
	Db   Db
}

func OpenDB() (*sql.DB, error) {
	present := false
	username, present := os.LookupEnv("DB_USERNAME")
	if !present {
		return nil, errors.New("env variable DB_USERNAME missing")
	}
	password, present := os.LookupEnv("DB_PASSWORD")
	if !present {
		return nil, errors.New("env variable DB_PASSWORD missing")
	}
	password = url.QueryEscape(password) //escaping in case of weird password
	host, present := os.LookupEnv("DB_HOST")
	if !present {
		return nil, errors.New("env variable DB_HOST missing")
	}
	name, present := os.LookupEnv("DB_NAME")
	if !present {
		return nil, errors.New("env variable DB_NAME missing")
	}
	ssl, present := os.LookupEnv("DB_SSL")
	if !present {
		return nil, errors.New("env variable DB_SSL missing")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=%s", username, password, host, name, ssl)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)

	duration, err := time.ParseDuration("800ms")
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewPostgresRepo(db *sql.DB) (*PostgresRepo, error) {
	return &PostgresRepo{
		Users: UsersModel{DB: db},
	}, nil
}

func (p PostgresRepo) GetAUser(username string) (*repo.User, error) {
	query := `
		SELECT 
			name
		FROM
			repo.user
		WHERE username = $1;
		`

	var user repo.User

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Second)
	defer cancel()

	rowContext := p.Users.DB.QueryRowContext(ctx, query, username)

	err := rowContext.Scan(
		&user.Name,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}

	return &user, nil
}
