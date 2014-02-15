package main

import (
	"os"
	"path"
	"fmt"
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/msbranco/goconfig"
	"git.notmuchmail.org/git/notmuch.git/bindings/go/src/notmuch"
)

type ResponseData struct {
	Content string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/count", CountHandler)
	//r.HandleFunc("/search", SearchHandler)
	//r.HandleFunc("/tag", TagHandler)
	// show, reply, insert
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

type CountData struct {
	Query string
	Count uint
}

func CountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-type", "application/json")

	//func OpenDatabase(path string, mode DatabaseMode) (*Database, Status) {
	db, status := get_notmuch_db()
	if status != notmuch.STATUS_SUCCESS {
		http.Error(w, "notmuch Oops", http.StatusInternalServerError)
	}
	query := db.CreateQuery("*")
	msgCount := query.CountMessages()

	data := CountData{"*", msgCount}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, "json Oops", http.StatusInternalServerError)
	}
	jsonText := string(jsonBytes[:])
	fmt.Fprintf(w, jsonText)
}

func get_notmuch_db() (*notmuch.Database, notmuch.Status) {
	// honor NOTMUCH_CONFIG
	home := os.Getenv("NOTMUCH_CONFIG")
	if home == "" {
		home = os.Getenv("HOME")
	}

	cfg, err := goconfig.ReadConfigFile(path.Join(home, ".notmuch-config"))
	if err != nil {
		//log.Fatalf("error loading config file:", err)
		return nil, notmuch.STATUS_FILE_ERROR
	}

	db_path, _ := cfg.GetString("database", "path")

	return notmuch.OpenDatabase(db_path, notmuch.DATABASE_MODE_READ_ONLY)
}
