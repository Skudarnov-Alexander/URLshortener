package methods

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"

	js "github.com/Skudarnov-Alexander/URLshortener/json"
	m "github.com/Skudarnov-Alexander/URLshortener/withdb"
)

var Tpl *template.Template

func LoadTpl() {
	Tpl = template.Must(template.ParseFiles("views/home.html", "views/applyProcess.html"))
}

// GetHome открывает главную страницу
func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Tpl.ExecuteTemplate(w, "home.html", nil)
	}
}

// PostLongURLform отправляет сдлинную ссылку в сервис через форму на фронте
func PostLongURLform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed! Only POST method is supported.", http.StatusMethodNotAllowed)
		return

	}

	// имитация получения JSON с фронта: из запроса парсим r.Form c маршалим в JSON
	jsonData, err := js.JSONfromFrontend(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Printf("JSON data:\n%s\n", jsonData)

	// валидация данных в JSON (URL и дней действия ссылки)
	data, err := js.JSONValid(jsonData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	k, err := m.InsertData(data, m.InternalDB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Tpl.ExecuteTemplate(w, "applyProcess.html", m.InternalDB[k]) //формат
}

// GetLongURLform получает оригинальную длинную ссылку -- > редиректит сразу на нее
func GetLongURLform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	path := strings.Split(r.URL.Path, "/")
	k := path[len(path)-1]

	if _, ok := m.InternalDB[k]; !ok {
		http.Error(w, "Sorry, this short URL is not exist", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, m.InternalDB[k].LongURL, http.StatusMovedPermanently)

}


//// API без фронта

// curl -X POST -H "Content-Type: application/json" \
// -d '{"longurl": "www.google.com", "expiredIn": "5"}' localhost:8080


func PostLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed! Only POST method is supported.", http.StatusMethodNotAllowed)
		return

	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	

	w.Write([]byte("\n***Data parsing from request.body - DONE\n"))
	w.Write(body)
	w.Write([]byte("\n"))

	
	

	// валидация данных в JSON (URL и дней действия ссылки)
	data, err := js.JSONValid(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("\n***Validation - DONE\n"))
	fmt.Println("\nData after validation\n", data)

	k, err := m.InsertData(data, m.InternalDB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = k
	w.Write([]byte("\n***Inserting to Database - DONE\n"))



	JSONfromDB, err := json.Marshal(m.InternalDB[k]) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write([]byte("\n***Marshalling from Database - DONE\n"))

	w.Write(JSONfromDB)
	
	
	w.Write([]byte("\nBye bye!\n"))
}

/*
// getShort обрабатывает запросы на получение короткой ссылки
func GetLongURLform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		// Код 405
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	path := strings.Split(r.URL.Path, "/")
	k := path[len(path)-1]

	if _, ok := m.InternalDB[k]; !ok {
		http.Error(w, "Sorry, this short URL is not exist", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, m.InternalDB[k].LongURL, http.StatusMovedPermanently)

} */


