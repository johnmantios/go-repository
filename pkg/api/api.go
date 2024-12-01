package api

import (
	log "github.com/sirupsen/logrus"
	"johnmantios.com/go-repository/pkg/service"
	"sync"
)

type GreetingUserAPI struct {
	Svc    service.GreetingUserService
	Env    string
	Wg     sync.WaitGroup
	Logger *log.Logger
}
