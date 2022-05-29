package models

import (
	"fmt"
	"strings"
	u "github.com/Skudarnov-Alexander/URLshortener/internal/url"
)

type DataBase map[string]u.UrlInfo

var InternalDB DataBase

func InitInternalDB() {
	InternalDB = make(DataBase)
}

//TODO переписать через пакет bcrypt, ч
		//тобы одинаковая ссылка от разных юзеров давала одинаковый хэш.
		//Хотя вероятность повтора сейчас очень мала.
		//Пока не отрабатывает сценарий оригинальный URL - оригинальный хэш.

// Парсинг запроса и сохранение данных в базу (мапу)
func (db DataBase) Insert(urlItem *u.UrlInfo) (err error) {
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


