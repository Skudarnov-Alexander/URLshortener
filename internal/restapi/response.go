package api

import (
	"encoding/json"
	"net/http"

	u "github.com/Skudarnov-Alexander/URLshortener/internal/url"
	l "github.com/Skudarnov-Alexander/URLshortener/internal/log"
	
)

type APIResponse struct {
	Code    int32  		`json:"code,omitempty"`
	Message string 		`json:"message,omitempty"`
	Data	u.UrlInfo 	`json:"data,omitempty"`
}

type APIResponseError struct {
	Code    int32  		`json:"code,omitempty"`
	Message string 		`json:"message,omitempty"`
}

func ResponseBadRequest(msg string, w http.ResponseWriter) {
	res := APIResponseError {
		Code: http.StatusBadRequest,
		Message: msg,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
	l.Logger.Print(msg)
}

func ResponseNotFound(msg string, w http.ResponseWriter) {
	res := APIResponseError {
		Code: http.StatusBadRequest,
		Message: msg,
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(res)
	l.Logger.Print(msg)
}

func ResponseMethodNotAllowed(msg string, w http.ResponseWriter) {
	res := APIResponseError {
		Code: http.StatusMethodNotAllowed,
		Message: msg,
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(res)
	l.Logger.Print(msg)
}




