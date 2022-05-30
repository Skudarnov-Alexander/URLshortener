package shorturl

import (
	"math/rand"
	"time"
)

//TODO переписать через пакет bcrypt

// makeShortURL генерирует случайную последовательность из 10 символов.
// Обязательно есть хотя бы одна цифра и буква
func MakeShortURL() string {
	rand.Seed(time.Now().UnixNano())

	digits := "0123456789_"
	letters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	all := digits + letters

	length := 10

	b := make([]byte, length)
	b[0] = digits[rand.Intn(len(digits))]
	b[1] = letters[rand.Intn(len(letters))]
	for i := 2; i < length; i++ {
		b[i] = all[rand.Intn(len(all))]
	}
	rand.Shuffle(len(b), func(i, j int) {
		b[i], b[j] = b[j], b[i]
	})

	return string(b)
}
