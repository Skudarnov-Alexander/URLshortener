package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	LongURL      string
	ExpDate      string
	CreationDate string
	ShortURL     string
}

type DataBase map[string]Data

var tpl *template.Template
var db DataBase

func init() {
	tpl = template.Must(template.ParseFiles("templates/home.html", "templates/applyProcess.html"))

}

func main() {
	db = make(DataBase) // Аллокация памяти мапы

	mux := http.NewServeMux()
	mux.HandleFunc("/", createShortURL)
	mux.HandleFunc("/short/", getShort) 
	http.ListenAndServe(":8080", mux)

}

// makeShortURL генерирует случайную последовательность из 10 символов
// Обязательно есть хотя бы одна цифра и буква
func makeShortURL() string {
	rand.Seed(time.Now().UnixNano())

	digits := "0123456789_"
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	all := digits + letters

	length := 10

	b := make([]byte, length)
	b[0] = digits[rand.Intn(len(digits))]
	b[1] = letters[rand.Intn(len(letters))]
	for i := 2; i < length; i++ {
		b[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})

	return string(b)
}

// Парсинг данных из формы, генерация случайной ссылки, упаковка данных в Data struct
func makeData(r *http.Request) (Data, string) {
	hash := makeShortURL()
	shortURL := "localhost:8080/short/" + hash
	longURL := r.FormValue("longURL")
	ExpDays := r.FormValue("ExpDays")
	timeCreation := time.Now()
	timeCreationS := timeCreation.Format("2006/01/02") //TODO вынести форматирование

	ExpDaysInt, err := strconv.Atoi(ExpDays)
	if err != nil {
		log.Print(err)
	}

	ExpDate := timeCreation.AddDate(0, 0, ExpDaysInt)
	ExpDateS := ExpDate.Format("2006/01/02")

	data := Data{
		LongURL:      longURL,
		ExpDate:      ExpDateS,
		CreationDate: timeCreationS,
		ShortURL:     shortURL,
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

		// TODO: здесь нужен мьютекс или мапу из sync https://pkg.go.dev/sync#Map
		db[k] = data

		fmt.Println(db) 

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
