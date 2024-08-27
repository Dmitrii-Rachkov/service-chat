package handler

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"

	"service-chat/internal/db/entity"
)

const (
	StatusOK    = "OK"
	StatusError = "Error"
)

type Response struct {
	Status       string                `json:"status"`
	Error        string                `json:"error,omitempty"`
	Message      string                `json:"message,omitempty"`
	MessagesList []entity.Message      `json:"messages_list,omitempty"`
	ChatsList    []entity.Chat         `json:"chats_list,omitempty"`
	DelChatsList []entity.DeletedChats `json:"del_chats_list,omitempty"`
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
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is a required field", err.Field()))
		case "max":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s cannot exceed 20 characters", err.Field()))
		case "min":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s must contain at least 6 characters", err.Field()))
		case "containsany":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s must contain Latin letters and Arabic numerals, as well as the symbols @#$&*()", err.Field()))
		case "excludesall":
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s must not contain symbols !@#$&*()?", err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("Field %s is not valid", err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error:  strings.Join(errMsgs, ", "),
	}
}
