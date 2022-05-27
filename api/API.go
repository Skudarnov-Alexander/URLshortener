package methods

import (
	//"encoding/json"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

	

	log.Printf("\n***Data parsing from request.body - DONE\nBody: %s\n", string(body))
	
	// валидация данных в JSON (URL и дней действия ссылки)
	data, err := js.JSONValid(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("\n***Validation - DONE\nData: %v\n", data)

	k, err := m.InsertData(data, m.InternalDB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("***Inserting to Database - DONE\n")

	JSONfromDB, err := json.Marshal(m.InternalDB[k]) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("***Marshalling from Database - DONE\n")

	w.Write(JSONfromDB)
	
}


// getLongURL обрабатывает запросы на получение короткой ссылки
func GetLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	// читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	data := js.LongURLget{}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("some JSON-marshaller error")
		return
	}

	// валидация ссылки в теле запроса
	_, err = url.ParseRequestURI(data.ShortURL)
	if err != nil {
		log.Println("invalid url")
		return
	}

	path := strings.Split(data.ShortURL, "/")
	k := path[len(path)-1]


	if _, ok := m.InternalDB[k]; !ok {
		http.Error(w, "Sorry, this short URL is not exist", http.StatusNotFound)
		return
	}

	JSONfromDB, err := json.Marshal(m.InternalDB[k].LongURL) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("***Marshalling from Database - DONE\n")

	w.Write(JSONfromDB)
} 


