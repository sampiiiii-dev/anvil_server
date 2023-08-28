package common

import (
	"encoding/json"
	"log"
	"net/http"
)

// CustomResponse formats the HTTP response in a custom way.
func CustomResponse(w http.ResponseWriter, status int, details interface{}, extraInfo interface{}) {
	if w == nil {
		log.Fatal("ResponseWriter is nil")
		return
	}

	response := map[string]interface{}{
		"status": "error",
		"payload": map[string]interface{}{
			"data":  nil,
			"error": nil,
		},
		"meta": map[string]interface{}{
			"message": nil,
		},
	}

	if status >= 200 && status < 300 {
		response["status"] = "success"
		response["payload"].(map[string]interface{})["data"] = details
		response["meta"].(map[string]interface{})["message"] = extraInfo
	} else {
		response["payload"].(map[string]interface{})["error"] = map[string]interface{}{
			"code":    status,
			"message": details,
		}
		response["meta"].(map[string]interface{})["message"] = extraInfo
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("JSON encoding failed: %s", err)

		// Send a standard API response format for Internal Server Error
		w.WriteHeader(http.StatusInternalServerError)
		internalServerErrorResponse := map[string]interface{}{
			"status": "error",
			"payload": map[string]interface{}{
				"data":  nil,
				"error": map[string]interface{}{"code": http.StatusInternalServerError, "message": "Internal Server Error"},
			},
			"meta": map[string]interface{}{"message": "An internal error occurred"},
		}
		err := json.NewEncoder(w).Encode(internalServerErrorResponse)
		if err != nil {
			panic(err)
		}
	}
}
