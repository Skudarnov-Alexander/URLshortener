package api

import (
	"encoding/json"
	"net/http"

	u "github.com/Skudarnov-Alexander/URLshortener/internal/url"
)

type APIResponse struct {
	Code    int32  		`json:"code"`
	Type    string 		`json:"type"`
	Message string 		`json:"message"`
	Data	u.UrlInfo 	`json:"data"`
}

func ResponseBadRequest(msg string, w http.ResponseWriter) {
	res := APIResponse {
		Code: http.StatusBadRequest,
		Message: msg,
	}

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(res)
}

func ResponseMethodNotAllowed(msg string, w http.ResponseWriter) {
	res := APIResponse {
		Code: http.StatusMethodNotAllowed,
		Message: msg,
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(res)
}




