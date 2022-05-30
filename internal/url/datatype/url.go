package datatype

import "time"

type UrlInfo struct {
	LongURL   string    `json:"longUrl"`
	ExpiredAt time.Time `json:"expiredAt"`
	CreatedAt time.Time `json:"createdAt"`
	ShortURL  string    `json:"shortURL"`
	ExpiredIn int       `json:"expiredIn"`
	// TODO добавить количество использований и флаг, что ссылка протухла
}
