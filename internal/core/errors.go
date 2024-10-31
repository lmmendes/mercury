package core

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type APIError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type ValidationError struct {
	Field  string `json:"field"`
	Rule   string `json:"rule"`
	Value  string `json:"value"`
	Reason string `json:"reason"`
}

func (c *Core) HandleError(err error, code int) error {
	if code >= 500 {
		c.Logger.ErrorWithStack(err)
	} else {
		c.Logger.Error("HTTP %d: %v", code, err)
	}

	// Handle validation errors
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		details := make([]ValidationError, 0)
		for _, err := range validationErrors {
			details = append(details, ValidationError{
				Field:  err.Field(),
				Rule:   err.Tag(),
				Value:  fmt.Sprintf("%v", err.Value()),
				Reason: err.Error(),
			})
		}
		return echo.NewHTTPError(http.StatusBadRequest, APIError{
			Code:    http.StatusBadRequest,
			Message: "Validation failed",
			Details: details,
		})
	}

	// Handle not found errors
	if err == ErrNotFound {
		return echo.NewHTTPError(http.StatusNotFound, APIError{
			Code:    http.StatusNotFound,
			Message: err.Error(),
		})
	}

	return echo.NewHTTPError(code, APIError{
		Code:    code,
		Message: err.Error(),
	})
}

// Common errors
var (
	ErrNotFound     = fmt.Errorf("resource not found")
	ErrUnauthorized = fmt.Errorf("unauthorized")
	ErrForbidden    = fmt.Errorf("forbidden")
)
