package api

import (
	"net/http"
)

func (api *GreetingUserAPI) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	env := Envelope{
		"status": "available",
		"system_info": map[string]string{
			"message": "Ready to greet!",
		},
	}

	err := WriteJSON(w, http.StatusOK, env, nil)
	if err != nil {
		serverErrorResponse(w, r, err, api.Logger)
	}
}
