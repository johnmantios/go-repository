package api

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (api *GreetingUserAPI) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(api.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(api.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", api.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/greet/:username", api.greetHandler)

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return api.logResponse(api.RecoverPanic(api.EnableCORS(api.logRequest(api.SecureHeaders(router)))))
}
