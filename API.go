package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	js "github.com/Skudarnov-Alexander/URLshortener/json"
)

// createShortURL обрабатывает запросы на добавление новой ссылки
func getHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "home.html", nil)
	}
}

func postLongURLform(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		// имитация получения JSON с фронта: из запроса парсим r.Form c маршалим в JSON
		jsonData, err := js.JSONfromFrontend(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Printf("JSON data:\n%s\n", jsonData)

		// валидация данных в JSON (URL и дней действия ссылки)
		data, err := js.JSONValid(jsonData)
		if err!= nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		k, err := InsertData(data, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		tpl.ExecuteTemplate(w, "applyProcess.html", db[k]) //формат

	} else {
		io.WriteString(w, "Sorry, only GET and POST methods are supported.\n")
	}
}

// getShort обрабатывает запросы на получение короткой ссылки
func getLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// Header
		w.Header().Set("Allow", http.MethodGet)

		// Код 405
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	path := strings.Split(r.URL.Path, "/")
	k := path[len(path)-1]

	if _, ok := db[k]; !ok {
		//io.WriteString(w, "Sorry, this short URL is not exist\n")
		http.Error(w, "Sorry, this short URL is not exist", http.StatusNotFound)
		return
	}
	url := "https://" + db[k].LongURL

	http.Redirect(w, r, url, http.StatusMovedPermanently)

}
