package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// структура, под которую будем принимать JSON
type Message struct {
	Seria     int    `json:"seria"`
	Nomer     int    `json:"nomer"`
	Umia      string `json:"umia"`
	Famlia    string `json:"famlia"`
	Otchestvo string `json:"otchestvo"`
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	// читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// парсим JSON
	var msg Message
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	// создаём папку, если её нет
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	// формируем имя файла по времени
	filename := fmt.Sprintf("data/%d.json", msg.Seria)

	// записываем JSON как есть
	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		http.Error(w, "Ошибка записи", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Файл сохранён: %s\n", filename)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET", http.StatusMethodNotAllowed)
		return
	}

	file := r.URL.Query().Get("seria")
	if file == "" {
		http.Error(w, "Не указан параметр seria", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("data/%s", file)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	var ms Message
	if err := json.Unmarshal(data, &ms); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusInternalServerError)
		return
	}

	resp := struct {
		Umia      string `json:"Имя"`
		Famlia    string `json:"Фамилия"`
		Otchestvo string `json:"Отчество"`
	}{
		Umia:      ms.Umia,
		Famlia:    ms.Famlia,
		Otchestvo: ms.Otchestvo,
	}

	// выставляем правильный Content-Type
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

type Messagev2 struct {
	Series     int    `json:"series"`
	Number     int    `json:"Number"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	MiddleName string `json:"middleName"`
}

func saveHandlerv2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	// читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// парсим JSON
	var msg Messagev2
	if err := json.Unmarshal(body, &msg); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusBadRequest)
		return
	}

	// создаём папку, если её нет
	if _, err := os.Stat("datav2"); os.IsNotExist(err) {
		os.Mkdir("datav2", 0755)
	}

	// формируем имя файла по времени
	filename := fmt.Sprintf("datav2/%d.json", msg.Series)

	// записываем JSON как есть
	if err := ioutil.WriteFile(filename, body, 0644); err != nil {
		http.Error(w, "Ошибка записи", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Файл сохранён: %s\n", filename)
}

func getHandlerv2(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET", http.StatusMethodNotAllowed)
		return
	}

	file := r.URL.Query().Get("series")
	if file == "" {
		http.Error(w, "Не указан параметр series", http.StatusBadRequest)
		return
	}

	path := fmt.Sprintf("datav2/%s", file)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	var mss Messagev2
	if err := json.Unmarshal(data, &mss); err != nil {
		http.Error(w, "Ошибка JSON", http.StatusInternalServerError)
		return
	}

	respp := struct {
		FirstName  string `json:"FirstName"`
		LastName   string `json:"LastName"`
		MiddleName string `json:"MiddleName"`
	}{
		FirstName:  mss.FirstName,
		LastName:   mss.LastName,
		MiddleName: mss.MiddleName,
	}

	// выставляем правильный Content-Type
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respp)
}

func main() {
	http.HandleFunc("/save", saveHandler)
	http.HandleFunc("/get", getHandler)

	http.HandleFunc("/v2/save", saveHandlerv2)
	http.HandleFunc("/v2/get", getHandlerv2)

	fmt.Println("Сервер запущен на :8080")
	http.ListenAndServe(":8080", nil)
}
