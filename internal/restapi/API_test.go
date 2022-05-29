package api

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	//js "github.com/Skudarnov-Alexander/URLshortener/json"
)



/*
func PostLongURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed! Only POST method is supported.", http.StatusMethodNotAllowed)
		return

	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Unable to parse request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()



	log.Printf("\n***Data parsing from request.body - DONE\nBody: %s\n", string(body))

	// валидация данных в JSON (URL и дней действия ссылки)
	data, err := js.JSONValid(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("\n***Validation - DONE\nData: %v\n", data)

	k, err := m.InsertData(data, m.InternalDB)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	log.Printf("***Inserting to Database - DONE\n")

	JSONfromDB, err := json.Marshal(m.InternalDB[k])

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("***Marshalling from Database - DONE\n")

	w.Write(JSONfromDB)

}

*/
func TestPostLongURL(t *testing.T) {
	// Создаем запрос с указанием нашего хендлера. Так как мы тестируем GET-эндпоинт
	// то нам не нужно передавать тело, поэтому третьим аргументом передаем nil


	testCase := struct {
		LongURL   string
		ExpiredIn int 
	}{
		LongURL:   "http://www.google.com",
		ExpiredIn: 5,

	}
	_ = testCase

	str := `{"longurl": "http://www.google.com", "expiredIn": 5}`

	b := bytes.NewBufferString(str)

	
	req, err := http.NewRequest("POST", "localhost:8080/", b)
	if err != nil {
		t.Fatal(err)
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(PostLongURL)

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `Parsed query-param with key "name": John`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
