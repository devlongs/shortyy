package main

import (
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Mapper struct {
	Mapping map[string]string
	Lock sync.Mutex
}

var urlMapper Mapper

func init() {
    urlMapper.Mapping = make(map[string]string)
}

func main() {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Server running..."))
	})

	http.ListenAndServe(":8080", r)
}

