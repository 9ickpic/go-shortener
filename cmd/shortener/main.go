package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
)

/*
	type Subj struct {
		Product string `json:"name"`
		Price   int    `json:"price"`
	}

	func JSONHandler(w http.ResponseWriter, req *http.Request) {
		// собираем данные
		subj := Subj{"Milk", 50}
		// code in JSON
		resp, err := json.Marshal(subj)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		// set header Content-Type
		// from give info client, code in JSON
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}

type MyHandler struct{}

// learn ya go-dev создание сервера HTTP, 3/10

	func mainPage(res http.ResponseWriter, req *http.Request) {
		body := fmt.Sprintf("Method: %s\r\n", req.Method)
		body += "Header ===============\r\n"
		for k, v := range req.Header {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}
		body += `Query parameters ===============\r\n`
		if err := req.ParseForm(); err != nil {
			res.Write([]byte(err.Error()))
			return
		}
		for k, v := range req.Form {
			body += fmt.Sprintf("%s: %v\r\n", k, v)
		}
		res.Write([]byte(body))
	}

	func apiPage(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Это страница /api."))
	}
*/
var urls map[string]string

func main() {
	urls = make(map[string]string)

	mux := http.NewServeMux()
	mux.HandleFunc(`/`, shortenURL)
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(shortenURL))
}

func shortenURL(res http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodPost {
		responseData, err := io.ReadAll(req.Body)

		if err != nil {
			http.Error(res, fmt.Sprintf("cannot request body: %s", err), http.StatusBadRequest)
			return
		}

		if string(responseData) == "" {
			http.Error(res, "empty POST request body", http.StatusBadRequest)
			return
		}

		url := string(responseData)
		if url == "" {
			http.Error(res, "empty url", http.StatusBadRequest)
			return
		}
		id := generateID()
		urls[id] = url
		response := fmt.Sprintf("http://localhost:8080/%s", id)
		res.Header().Set("Content-type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		_, err = res.Write([]byte(response))
		if err != nil {
			return
		}
	} else if req.Method == http.MethodGet {
		id := req.URL.Path[1:]
		url, ok := urls[id]
		if !ok {
			http.Error(res, "invalid url", http.StatusBadRequest)
		}
		res.Header().Set("Location", url)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		res.WriteHeader(http.StatusBadRequest)
	}
}

func generateID() string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, 8)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
