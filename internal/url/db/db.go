package db

import (
	"fmt"
	"github.com/Skudarnov-Alexander/URLshortener/internal/url/datatype"
	"strings"
)

type DataBase map[string]datatype.UrlInfo

func New() *DataBase {
	db := make(DataBase)
	return &db
}

// Cохранение данных в базу (мапу)
func (db DataBase) Insert(urlItem *datatype.UrlInfo) (err error) {
	path := strings.Split(urlItem.ShortURL, "/")
	k := path[len(path)-1]

	if _, ok := db[k]; ok {
		err = fmt.Errorf("this short URL is already exist")
		return
	}

	// TODO: здесь нужен мьютекс или мапу из sync https://pkg.go.dev/sync#Map
	db[k] = *urlItem
	return
}
