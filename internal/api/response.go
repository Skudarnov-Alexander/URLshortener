package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Skudarnov-Alexander/URLshortener/internal/url/datatype"
)

type APIResponse struct {
	Code    int32            `json:"code,omitempty"`
	Message string           `json:"message,omitempty"`
	Data    datatype.UrlInfo `json:"data,omitempty"`
}

type APIResponseError struct {
	Code    int32  `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func ResponseBadRequest(l *log.Logger, msg string, w http.ResponseWriter) {
	res := APIResponseError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
	l.Print(msg)
}

func ResponseNotFound(l *log.Logger, msg string, w http.ResponseWriter) {
	res := APIResponseError{
		Code:    http.StatusBadRequest,
		Message: msg,
	}
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(res)
	l.Print(msg)
}

func ResponseMethodNotAllowed(l *log.Logger, msg string, w http.ResponseWriter) {
	res := APIResponseError{
		Code:    http.StatusMethodNotAllowed,
		Message: msg,
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(res)
	l.Print(msg)
}
