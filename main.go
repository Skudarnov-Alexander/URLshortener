package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	api "github.com/Skudarnov-Alexander/URLshortener/api"
	m "github.com/Skudarnov-Alexander/URLshortener/withdb"
)


var Tpl *template.Template
var Db m.DataBase
var MyLogger *log.Logger

func init() {
	api.LoadTpl()		// Загрузка HTML
	m.InitInternalDB()  // Аллокация памяти мапы для хранения ссылок
}

func main() {

	MyLogger = log.New(os.Stdout, "log", 19)
	



	mux := http.NewServeMux()
	mux.HandleFunc("/form/", api.GetHome)
	//mux.HandleFunc("/form/ok", api.PostLongURLform)
	mux.HandleFunc("/short/", api.GetLongURLform)
	mux.HandleFunc("/", api.PostLongURL)
	mux.HandleFunc("/sh/", api.GetLongURL)
	MyLogger.Printf("Server starting...")
	http.ListenAndServe(":8080", mux)
	

}
