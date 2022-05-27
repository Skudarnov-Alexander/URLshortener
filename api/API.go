package methods

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	js "github.com/Skudarnov-Alexander/URLshortener/json"
	m "github.com/Skudarnov-Alexander/URLshortener/withdb"
)

var Tpl *template.Template

func LoadTpl() {
	Tpl = template.Must(template.ParseFiles("views/home.html", "views/applyProcess.html"))
}

// createShortURL обрабатывает запросы на добавление новой ссылки
func GetHome(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		Tpl.ExecuteTemplate(w, "home.html", nil)
	}
}

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

// getShort обрабатывает запросы на получение короткой ссылки
func GetLongURL(w http.ResponseWriter, r *http.Request) {
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

}
