package api

import (
	//"encoding/json"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	datatype "github.com/Skudarnov-Alexander/URLshortener/internal/url"
	js "github.com/Skudarnov-Alexander/URLshortener/json"
	sh "github.com/Skudarnov-Alexander/URLshortener/utils"
	m "github.com/Skudarnov-Alexander/URLshortener/withdb"
)

var Tpl *template.Template

func LoadTpl() {
	Tpl = template.Must(template.ParseFiles("views/home.html", "views/applyProcess.html"))
}

// GetHome открывает главную страницу
func GetHome(w http.ResponseWriter, r *http.Request) {
	//MyLogger.Println("GET GET GET")
	if r.Method == http.MethodGet {
		Tpl.ExecuteTemplate(w, "home.html", nil)
	}
}
/*
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
*/
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
		ResponseMethodNotAllowed("Only POST method is supported.", w)
		return
	}

	URLItem := new(datatype.UrlInfo)

	err := json.NewDecoder(r.Body).Decode(URLItem)
	if err != nil {
		ResponseBadRequest("Parsing request body error", w)
	}

	// JSON data validation
	if URLItem.LongURL == "" {
		ResponseBadRequest("Empty URL", w)
		return
	}

	if _, err = url.ParseRequestURI(URLItem.LongURL); err != nil {
		ResponseBadRequest("Incorrect URL", w)
		return
	}

	if URLItem.ExpiredIn < 0 {
		ResponseBadRequest("Days have negative value", w)
		return
	}

	// Preparing additional data
	k := sh.MakeShortURL()
	URLItem.ShortURL = "localhost:8080/short/" + k
	URLItem.CreatedAt = time.Now()
	if URLItem.ExpiredIn == 0 {
		URLItem.ExpiredIn = 365
	}

	URLItem.ExpiredAt = URLItem.CreatedAt.AddDate(0, 0, URLItem.ExpiredIn)

	// add item to Database
	err = m.InternalDB.Insert(URLItem)
	if err != nil {
		ResponseBadRequest("This short URL is already exist", w)
		return
	}
	
	resp := APIResponse{
		Code:	http.StatusCreated,  				
		Data:	*URLItem,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
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

	w.Write(JSONfromDB)
} 


