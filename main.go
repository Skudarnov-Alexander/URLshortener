package main

import (
	api "github.com/Skudarnov-Alexander/URLshortener/api"
	m "github.com/Skudarnov-Alexander/URLshortener/withdb"
	"html/template"
	"net/http"
)

var Tpl *template.Template
var Db m.DataBase

func init() {
	api.LoadTpl()
	m.InitInternalDB() // Аллокация памяти мапы для хранения ссылок
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/form/", api.GetHome)
	mux.HandleFunc("/form/ok", api.PostLongURLform)
	mux.HandleFunc("/short/", api.GetLongURL)
	http.ListenAndServe(":8080", mux)

}
