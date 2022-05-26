package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template
var db DataBase

func init() {
	tpl = template.Must(template.ParseFiles("views/home.html", "views/applyProcess.html"))

}

func main() {
	db = make(DataBase) // Аллокация памяти мапы для хранения ссылок

	mux := http.NewServeMux()
	mux.HandleFunc("/form/", getHome)
	mux.HandleFunc("/form/ok", postLongURL)
	mux.HandleFunc("/form/short/", getLongURL)
	http.ListenAndServe(":8080", mux)

}
