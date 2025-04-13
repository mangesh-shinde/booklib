package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	statusOK    = "OK"
	statusError = "Error"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func SendError(w http.ResponseWriter, statusCode int, errorMessage error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := Response{Status: statusError, Error: errorMessage.Error()}
	json.NewEncoder(w).Encode(resp)
}

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func ValidateErrors(errs validator.ValidationErrors) Response {
	var errMessages []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMessages = append(errMessages, fmt.Sprintf("field %s is required", err.Field()))
		default:
			errMessages = append(errMessages, fmt.Sprintf("field %s is invalid", err.Field()))
		}
	}

	return Response{
		Status: statusError,
		Error:  strings.Join(errMessages, ", "),
	}

}
