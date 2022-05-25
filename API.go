package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

// createShortURL обрабатывает запросы на добавление новой ссылки
func createShortURL(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		tpl.ExecuteTemplate(w, "home.html", nil)

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 500)
			// TODO: отдавать ответ клиенту
			return
		}

		k, err := InsertData(r, db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		//TODO форматирование time.Time

		jsonData, err := json.MarshalIndent(db[k], "", " ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Println("json:", jsonData)

		newJson := UrlData{}

		err = json.Unmarshal(jsonData, &newJson)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Println(newJson)

		tpl.ExecuteTemplate(w, "applyProcess.html", db[k])
	default:
		io.WriteString(w, "Sorry, only GET and POST methods are supported.\n")

	}
}

// getShort обрабатывает запросы на получение короткой ссылки
func getShort(w http.ResponseWriter, r *http.Request) {
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