package api

import (
	"fmt"
	"net/http"
)

func (api *GreetingUserAPI) greetHandler(w http.ResponseWriter, r *http.Request) {
	username, err := readUsernameParam(r)
	if err != nil {
		api.notFoundResponse(w, r)
		return
	}

	user, err := api.Svc.Greet(username)
	if err != nil {
		serverErrorResponse(w, r, err, api.Logger)
		return
	}

	headers := make(http.Header)
	headers.Set("ID", fmt.Sprintf("/v1/greet/%s", user.Name))

	err = WriteJSON(w, http.StatusOK, Envelope{"name": user.Name}, headers)
	if err != nil {
		serverErrorResponse(w, r, err, api.Logger)
		return
	}

}
