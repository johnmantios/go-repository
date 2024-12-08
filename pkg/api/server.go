package api

import (
	"context"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (api *GreetingUserAPI) Serve() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", 4000),
		Handler:      api.Routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		api.Logger.WithFields(log.Fields{
			"caught signal": s.String(),
		}).Info()

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		api.Logger.WithFields(log.Fields{
			"completing background tasks": s.String(),
			"addr":                        srv.Addr,
		}).Info()

		api.Wg.Wait()
		shutdownError <- nil
	}()

	api.Logger.WithFields(log.Fields{
		"env":  api.Env,
		"addr": srv.Addr,
	}).Info("Starting server...")

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	api.Logger.WithFields(log.Fields{
		"env":  api.Env,
		"addr": srv.Addr,
	}).Info("stopped server")

	return nil
}
