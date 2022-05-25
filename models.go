package main

import "time"

type Data struct {
	LongURL      string		 	`json:"longUrl"`
	ExpiredAt    time.Time		`json:"expiredAt"`
	CreatedAt 	 time.Time		`json:"createdAt"`
	ShortURL     string			`json:"shortURL"`
}

type DataBase map[string]Data