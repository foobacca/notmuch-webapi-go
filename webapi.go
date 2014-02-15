package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
)

type ResponseData struct {
	Content string
	//Content []byte
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	data := ResponseData{"hello world"}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "Oops", http.StatusInternalServerError)
	}
	jsonText := string(jsonBytes[:])
	fmt.Fprintf(w, jsonText)
}
