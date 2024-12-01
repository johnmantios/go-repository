package api

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func serverErrorResponse(w http.ResponseWriter, r *http.Request, err error, logger *log.Logger) {
	logger.Println(r, err)
	logger.WithFields(log.Fields{
		"request": r,
	}).Error(err)

	message := "the server encountered a problem and could not process your request"
	errorResponse(w, r, http.StatusInternalServerError, message, logger)
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, message any, logger *log.Logger) {
	env := Envelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		logger.WithFields(log.Fields{
			"request": r,
		}).Error(err)
		w.WriteHeader(500)
	}
}

func (api *GreetingUserAPI) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	errorResponse(w, r, http.StatusNotFound, message, api.Logger)
}

func (api *GreetingUserAPI) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	errorResponse(w, r, http.StatusMethodNotAllowed, message, api.Logger)
}

func badRequestResponse(w http.ResponseWriter, r *http.Request, err error, logger *log.Logger) {
	errorResponse(w, r, http.StatusBadRequest, err.Error(), logger)
}
