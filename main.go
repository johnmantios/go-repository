package main

import (
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"johnmantios.com/go-repository/pkg/api"
	repository "johnmantios.com/go-repository/pkg/repo/redis"
	"johnmantios.com/go-repository/pkg/service"
	"log"
)

func main() {

	db, err := repository.OpenDB()
	if err != nil {
		log.Panic(err)
	}

	repo, err := repository.NewRedisRepo(db)
	if err != nil {
		log.Panic(err)
	}

	svc := service.NewGreetingUserService(repo)

	var logger = logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		// DisableTimestamp: false, // Default: false
		// DisableSorting:   false, // sort the fields in the order of the keys
		PrettyPrint:     true, // for debugging
		TimestampFormat: "2006-01-02 15:04:05",
	})

	api := &api.GreetingUserAPI{
		Svc:    svc,
		Env:    "dev",
		Logger: logger,
	}

	err = api.Serve()
	if err != nil {
		log.Panic(err)
	}
}
