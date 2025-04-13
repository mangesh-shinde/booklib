package response

import (
	"encoding/json"
	"net/http"
)

func SendError(w http.ResponseWriter, statusCode int, errorMessage string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	resp := map[string]string{"error": errorMessage}
	json.NewEncoder(w).Encode(resp)
}
