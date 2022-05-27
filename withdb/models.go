package models

import (
	"fmt"
	"log"
	"time"

	js "github.com/Skudarnov-Alexander/URLshortener/json"
	sh "github.com/Skudarnov-Alexander/URLshortener/utils"
)

type UrlData struct {
	LongURL   string    `json:"longUrl"`
	ExpiredAt time.Time `json:"expiredAt"`
	CreatedAt time.Time `json:"createdAt"`
	ShortURL  string    `json:"shortURL"`
}

type DataBase map[string]UrlData

var InternalDB DataBase

func InitInternalDB() {
	InternalDB = make(DataBase)
}

// Парсинг запроса и сохранение данных в базу (мапу)
func InsertData(data js.LongURL, db DataBase) (key string, err error) {
	key = sh.MakeShortURL()
	shortURL := "localhost:8080/short/" + key
	createdAt := time.Now()

	//Если пользователь ввел 0 количество дней, то по умолчанию ссылка действительна 1 год.
	if data.ExpiredIn == 0 {
		data.ExpiredIn = 365
	}

	expiredAt := createdAt.AddDate(0, 0, data.ExpiredIn)

	// Запись в базу
	if _, ok := db[key]; ok {
		//TODO переписать через пакет bcrypt, ч
		//тобы одинаковая ссылка от разных юзеров давала одинаковый хэш.
		//Хотя вероятность повтора сейчас очень мала.
		//Пока не отрабатывает сценарий оригинальный URL - оригинальный хэш.
		log.Print("Повтор ключа. Перезапись ссылки")
		err = fmt.Errorf("повтор ключа. Перезапись ссылки")
		return
	}

	// TODO: здесь нужен мьютекс или мапу из sync https://pkg.go.dev/sync#Map
	db[key] = UrlData{
		LongURL:   data.LongURL,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
		ShortURL:  shortURL,
	}

	return
}

// curl -X POST -H "Content-Type: application/json" \
// -d '{"longurl": "www.google.com", "expiredIn": "5"}' localhost:8080
