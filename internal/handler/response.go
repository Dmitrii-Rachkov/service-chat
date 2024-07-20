package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type errorResponse struct {
	Message string `json:"message"`
}

type Response struct {
	Status  string `json:"status"`
	Error   string `json:"error,omitempty"`
	Message string `json:"message,omitempty"`
}

func OK(msg string) Response {
	return Response{
		Status:  StatusOK,
		Message: msg,
	}
}

func Error(msg string) Response {
	return Response{
		Status: StatusError,
		Error:  msg,
	}
}

func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is a required field", err.Field()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s cannot exceed 20 characters", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must contain at least 6 characters", err.Field()))
		case "containsany":
			errMsgs = append(errMsgs, fmt.Sprintf("field %s must contain Latin letters and Arabic numerals, as well as the symbols @#$&*()", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
