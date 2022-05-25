package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UrlData struct {
	LongURL   string    `json:"longUrl"`
	ExpiredAt time.Time `json:"expiredAt"`
	CreatedAt time.Time `json:"createdAt"`
	ShortURL  string    `json:"shortURL"`
}

type DataBase map[string]UrlData

// Парсинг запроса и сохранение данных в базу (мапу)
func InsertData(r *http.Request, db DataBase) (key string, err error) {
	key = MakeShortURL()
	shortURL := "localhost:8080/short/" + key
	longURL := r.FormValue("longURL") //TODO валидация на протокол http или https
	if longURL == "" {
		log.Print("Пустая ссылка")
		err = fmt.Errorf("пустая ссылка")
		return
	}
	expDays := r.FormValue("ExpDays")
	createdAt := time.Now()

	expDaysInt, err := strconv.Atoi(expDays)
	if err != nil {
		log.Print("Ошибка. Введите число!")
		return
	}

	expiredAt := createdAt.AddDate(0, 0, expDaysInt)

	if _, ok := db[key]; !ok {
		//TODO переписать через пакет bcrypt, чтобы одинаковая ссылка от разных юзеров давала одинаковый хэш.
		//Хотя вероятность повтора сейчас очень мала. Пока не отрабатывает сценарий оригинальный URL - оригинальный хэш.
		log.Print("Повтор ключа. Перезапись ссылки")
		err = fmt.Errorf("повтор ключа. Перезапись ссылки")
		return
	}
	// TODO: здесь нужен мьютекс или мапу из sync https://pkg.go.dev/sync#Map
	db[key] = UrlData{
		LongURL:   longURL,
		ExpiredAt: expiredAt,
		CreatedAt: createdAt,
		ShortURL:  shortURL,
	}

	return
}
