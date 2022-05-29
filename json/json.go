// тестовый пакет для работы с фронтом

package json

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type LongURL struct {
	LongURL   string `json:"longurl"`
	ExpiredIn int    `json:"expiredIn"`
}

type LongURLget struct {
	ShortURL   string `json:"shorturl"`
}

// JSONfromFrontend парсит значения от фронта из r.FormValue и маршалит в JSON (type LongURL)
func JSONfromFrontend(r *http.Request) (jsonData []byte, err error) {
	if err = r.ParseForm(); err != nil {
		log.Println("parse form error:", err)
		return
	}
	longURL := r.FormValue("longURL")
	expDays := r.FormValue("exp_days")

	expDaysInt, err := strconv.Atoi(expDays)
	if err != nil {
		log.Print("Incorrect Expired Days (not int):", err)
		//http.Error(w, "Incorrect Expired Days", http.StatusBadRequest)
		return
	}
	if expDaysInt < 0 {
		log.Print("Expired Days is Negative number:", err)
		err = fmt.Errorf("incorrect expired days (zero or negative)")
		//http.Error(w, "Incorrect Expired Days (less 1)", http.StatusBadRequest)
		return
	}

	data := LongURL{longURL, expDaysInt}

	jsonData, err = json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Println("some JSON-marshaller error")
		err = fmt.Errorf("some JSON-marshaller error")
		return
	}

	return
}

// TODO: валидация на протокол http или https
// todo если не ввел дату
// Валидатор и маршаллер данных с формы после отправления длинной ссылки
// todo valid empty
func JSONValid(jsonData []byte) (data LongURL, err error) {

	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		log.Println("some JSON-marshaller error")
		return
	}
	//TODO вывести на фронт
	if data.LongURL == "" {
		log.Print("Пустая ссылка")
		err = fmt.Errorf("пустая ссылка")
		return
	}

	_, err = url.ParseRequestURI(data.LongURL)
	if err != nil {
		log.Println("invalid url")
		return
	}

	if data.ExpiredIn < 0 {
		log.Print("Incorrect Expired Days (less 1):", err)
		err = fmt.Errorf("incorrect expired days (zero or negative)")
		return
	}

	return
}
