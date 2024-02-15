package main

import (
	filters "gohandler/pkg/filters"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	filters.AddFilterHandlers(mux)
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	mux.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok\n"))
	})
	log.Fatal(http.ListenAndServe("127.1:8080", mux))

}
