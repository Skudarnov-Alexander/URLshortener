package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)


var tpl *template.Template
var db DataBase

func init() {
	tpl = template.Must(template.ParseFiles("views/home.html", "views/applyProcess.html"))

}

func main() {
	db = make(DataBase) // Аллокация памяти мапы для хранения ссылок

	mux := http.NewServeMux()
	mux.HandleFunc("/", createShortURL)
	mux.HandleFunc("/short/", getShort) 
	http.ListenAndServe(":8080", mux)

}




// Парсинг данных из формы, генерация случайной ссылки, упаковка данных в Data struct
func makeData(r *http.Request) (Data, string) {
	hash := MakeShortURL()
	shortURL := "localhost:8080/short/" + hash
	longURL := r.FormValue("longURL")
	expDays := r.FormValue("ExpDays")
	createdAt := time.Now()
	

	expDaysInt, err := strconv.Atoi(expDays)
	if err != nil {
		log.Print("Ошибка. Введите число!", err)
	}

	expiredAt := createdAt.AddDate(0, 0, expDaysInt)
	

	data := Data{
		LongURL:        longURL,
		ExpiredAt:      expiredAt,
		CreatedAt:		createdAt,
		ShortURL:       shortURL,
	}

	return data, hash

}

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
			fmt.Println(err)
			// TODO: отдавать ответ клиенту
			return
		}
		fmt.Println(r.Form)

		data, k := makeData(r)

		fmt.Println(data)

		//timeCreationS := timeCreation.Format("2006/01/02") //TODO вынести форматирование
		//ExpDateS := ExpDate.Format("2006/01/02")

		// TODO: здесь нужен мьютекс или мапу из sync https://pkg.go.dev/sync#Map
		db[k] = data

		fmt.Println(db) 

		jsonData, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		fmt.Println("json:", jsonData)

		newJson := Data{}

		err = json.Unmarshal(jsonData, &newJson)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}


		fmt.Println(newJson)



		tpl.ExecuteTemplate(w, "applyProcess.html", data)
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

	fmt.Println("Redirect ", db[k].LongURL)

	http.Redirect(w, r, url, http.StatusMovedPermanently)

}
