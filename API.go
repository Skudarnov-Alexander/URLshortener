package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// createShortURL обрабатывает запросы на добавление новой ссылки
func getHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "home.html", nil)
	}
}

func postLongURL(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), 500)
			// TODO: отдавать ответ клиенту
			return
		}

		fmt.Println("db до вставки", db)

		k, err := InsertData(r, db)
		fmt.Println("db после вставки",db)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

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
		tpl.ExecuteTemplate(w, "applyProcess.html", db[k])

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
