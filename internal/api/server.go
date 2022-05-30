package api

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Skudarnov-Alexander/URLshortener/internal/logger"
	"github.com/Skudarnov-Alexander/URLshortener/internal/url/db"
)

type Config struct {
	RESTAPIport int
	RESTAPIhost string
	LogPath     string
}

// TODO разобраться с внешними пакетами по работе с конфигами
func NewConf(path string) (config *Config, err error) {
	config = new(Config)
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	//b := bytes.Buffer{}

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		txt := sc.Text()
		sf := strings.Fields(txt)
		switch sf[0] {
		case "PORT":
			config.RESTAPIport, _ = strconv.Atoi(sf[1])
		case "HOST":
			config.RESTAPIhost = sf[1]
		case "LOG_PATH":
			config.LogPath = sf[1]
		}
	}
	return
}

type RestAPI struct {
	Logger *logger.Logger
	Db     *db.DataBase
	Config *Config
	Mux    *http.ServeMux
}

func New(logger *logger.Logger, db *db.DataBase, cf *Config) *RestAPI {
	mux := http.NewServeMux()

	urlAPI := &RestAPI{
		Logger: logger,
		Db:     db,
		Config: cf,
		Mux:    mux,
	}

	urlAPI.Mux.HandleFunc("/", urlAPI.PostLongURL)
	//urlAPI.Mux.HandleFunc("/sh/", urlAPI.GetLongURL)
	return urlAPI

}

func (R *RestAPI) Start() {
	R.Logger.L.Print("Server is starting.....")
	port := ":" + strconv.Itoa(R.Config.RESTAPIport)
	http.ListenAndServe(port, R.Mux)
}
