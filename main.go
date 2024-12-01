package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"johnmantios.com/go-repository/pkg/api"
	repository "johnmantios.com/go-repository/pkg/repo/postgres"
	"johnmantios.com/go-repository/pkg/service"
	"os"
)

func main() {
	log.SetReportCaller(true)

	db, err := repository.OpenDB()
	if err != nil {
		log.Panic(err)
	}

	repo, err := repository.NewPostgresRepo(db)
	if err != nil {
		log.Panic(err)
	}

	svc := service.NewGreetingUserService(repo)

	logger := &log.Logger{
		Out: os.Stdout,
	}

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
