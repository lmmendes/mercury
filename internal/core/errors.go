package core

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Core) HandleError(w http.ResponseWriter, err error, code int) {
	// Log the error with stack trace for 5xx errors
	if code >= 500 {
		c.Logger.ErrorWithStack(err)
	} else {
		c.Logger.Error("HTTP %d: %v", code, err)
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{Code: code, Message: err.Error()})
}
