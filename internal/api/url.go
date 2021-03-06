package api

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Skudarnov-Alexander/URLshortener/internal/url/datatype"
	"github.com/Skudarnov-Alexander/URLshortener/internal/utils/shorturl"
)

//// API без фронта

// curl -X POST -H "Content-Type: application/json" \
// -d '{"longurl": "https://github.com", "expiredIn": "5"}' localhost:8080

func (u *RestAPI) PostLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		msg := "Only POST method is supported."
		ResponseMethodNotAllowed(u.Logger.L, msg, w)
		return
	}

	URLItem := new(datatype.UrlInfo)

	err := json.NewDecoder(r.Body).Decode(URLItem)
	if err != nil {
		msg := "Parsing request body error: " + err.Error()
		ResponseBadRequest(u.Logger.L, msg, w)
		return
	}

	// JSON data validation
	if URLItem.LongURL == "" {
		msg := "Empty URL"
		ResponseBadRequest(u.Logger.L, msg, w)
		return
	}

	if _, err = url.ParseRequestURI(URLItem.LongURL); err != nil {
		msg := "Incorrect URL: " + err.Error()
		ResponseBadRequest(u.Logger.L, msg, w)
		return
	}

	if URLItem.ExpiredIn < 0 {
		msg := "Days have negative value"
		ResponseBadRequest(u.Logger.L, msg, w)
		return
	}

	// Preparing additional data
	k := shorturl.MakeShortURL()
	host := u.Config.RESTAPIhost + ":" + strconv.Itoa(u.Config.RESTAPIport) + "/"
	URLItem.ShortURL = host + k
	URLItem.CreatedAt = time.Now()

	if URLItem.ExpiredIn == 0 {
		URLItem.ExpiredIn = 365
	}

	URLItem.ExpiredAt = URLItem.CreatedAt.AddDate(0, 0, URLItem.ExpiredIn)

	// add item to Database
	err = u.Db.Insert(URLItem)
	if err != nil {
		msg := "This short URL is already exist"
		ResponseBadRequest(u.Logger.L, msg, w)
		return
	}

	resp := APIResponse{
		Code:    http.StatusCreated,
		Data:    *URLItem,
		Message: "Short URL is created",
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

}

// curl -X GET -H "Content-Type: application/json" \
// -d '{"shorturl": "????"}' localhost:8080/sh/

// getLongURL обрабатывает запросы на получение короткой ссылки
/*
func GetLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	req := new(APIRequest)

	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		fmt.Println(err.Error())
		ResponseBadRequest("Parsing request body error:"+err.Error(), w)
		return
	}

	// валидация ссылки в теле запроса
	_, err = url.ParseRequestURI(req.ShortURL)
	if err != nil {
		ResponseBadRequest("invalid short url: "+err.Error(), w)
		return
	}

	path := strings.Split(req.ShortURL, "/")
	k := path[len(path)-1]

	if _, ok := InitDataBase[k]; !ok {
		msg := "Link is not found: " + k
		ResponseNotFound(msg, w)
		return
	}

	URLItem := new(datatype.UrlInfo)
	URLItem.LongURL = InitDataBase[k].LongURL

	resp := APIResponse{
		Code: http.StatusOK,
		Data: *URLItem,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

*/

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

	k, err := m.InsertData(data, InitDataBase)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	Tpl.ExecuteTemplate(w, "applyProcess.html", InitDataBase[k]) //формат
}

// GetLongURLform получает оригинальную длинную ссылку -- > редиректит сразу на нее
func GetLongURLform(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method is not allowed! Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	path := strings.Split(r.URL.Path, "/")
	k := path[len(path)-1]

	if _, ok := InitDataBase[k]; !ok {
		http.Error(w, "Sorry, this short URL is not exist", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, InitDataBase[k].LongURL, http.StatusMovedPermanently)

}

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

*/
