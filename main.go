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

func init() {
	api.LoadTpl()		// Загрузка HTML
	m.InitInternalDB()  // Аллокация памяти мапы для хранения ссылок
}

func main() {
	nf, err := os.Create("logs")
	if err != nil {
		log.Print("Error: creation log file", err)
		return
	}
	log.SetOutput(nf)
	log.SetPrefix("LOG:")
	log.SetFlags(19)
	defer nf.Close()



	mux := http.NewServeMux()
	mux.HandleFunc("/form/", api.GetHome)
	mux.HandleFunc("/form/ok", api.PostLongURLform)
	mux.HandleFunc("/short/", api.GetLongURLform)
	mux.HandleFunc("/", api.PostLongURL)
	//mux.HandleFunc("/sh/", api.GetLongURL)
	log.Printf("Server starting...")
	http.ListenAndServe(":8080", mux)
	

}
