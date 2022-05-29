package main

import (
	"html/template"
	"log"
	"net/http"

	api "github.com/Skudarnov-Alexander/URLshortener/internal/restapi"
	m "github.com/Skudarnov-Alexander/URLshortener/internal/url/withdb"
	l "github.com/Skudarnov-Alexander/URLshortener/internal/log"
)


var Tpl *template.Template
var Db m.DataBase
var MyLogger *log.Logger

func init() {
	//api.LoadTpl()		// Загрузка HTML
	m.InitInternalDB()  // Аллокация памяти мапы для хранения ссылок
}

func main() {
	m.InitInternalDB()
	l.InitLogger()
	



	mux := http.NewServeMux()
	//mux.HandleFunc("/form/", api.GetHome)
	//mux.HandleFunc("/form/ok", api.PostLongURLform)
	//mux.HandleFunc("/short/", api.GetLongURLform)
	mux.HandleFunc("/", api.PostLongURL)
	mux.HandleFunc("/sh/", api.GetLongURL)
	l.Logger.Print("Server is starting.....")
	http.ListenAndServe(":8080", mux)
	

}
