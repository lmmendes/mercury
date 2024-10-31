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
	c.Logger.Println(err)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(APIError{Code: code, Message: err.Error()})
}
