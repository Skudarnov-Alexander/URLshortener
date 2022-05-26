package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	sh "github.com/Skudarnov-Alexander/URLshortener/utils"
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
	key = sh.MakeShortURL()
	fmt.Println(key)
	shortURL := "localhost:8080/form/short/" + key
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

	v, ok := db[key]
	fmt.Printf("v: %v\nkey: %v\n", v, ok)

	if  ok {
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

type LongURL struct {
	LongURL   string `json:"longurl"`
	ExpiredIn string `json:"expiredIn"`
}

// curl -X POST -H "Content-Type: application/json"  -d '{"longurl": "www.google.com", "expiredIn": "5"}' localhost:8080

// TODO: Что вернет в JSON?
// TODO: Вывести ошибку в http.Error() в коде
// Валидатор и маршаллер данных с формы после отправления длинной ссылки
func JSONValid(r *http.Request) (jsonData []byte, err error) {
	longURL := r.FormValue("longURL") //TODO валидация на протокол http или https
	if longURL == "" {
		log.Print("Пустая ссылка")
		//http.Error(w, "Пустая ссылка", http.StatusBadRequest)
		err = fmt.Errorf("пустая ссылка")
		return 
	}

	//Если пользователь не ввел количество дней, то по умолчанию ссылка действительна 1 год.
	expDays := r.FormValue("ExpDays")
	
	if expDays == "" {
		expDays = "365"
	}

	expDaysInt, err := strconv.Atoi(expDays)
	if err != nil {
		log.Print("Incorrect Expired Days (not int):", err)
		//http.Error(w, "Incorrect Expired Days", http.StatusBadRequest)
		return 
	}
	if expDaysInt < 1 {
		log.Print("Incorrect Expired Days (less 1):", err)
		err = fmt.Errorf("incorrect expired days (zero or negative)")
		//http.Error(w, "Incorrect Expired Days (less 1)", http.StatusBadRequest)
		return 
	}

	data := LongURL{longURL, expDays}

	jsonData, err = json.MarshalIndent(data, "", " ")
		if err != nil {
			log.Print("JSON marshaler error")
			//http.Error(w, "JSON marshaler error", http.StatusBadRequest)
			return
		}
	
	return
}

