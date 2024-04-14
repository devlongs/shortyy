package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/lithammer/shortuuid/v4"
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

	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("Server running..."))
	})

	r.Post("/short-it", createShortURLHandler)
	r.Get("/short/{key}", redirectHandler)

	http.ListenAndServe(":8080", r)
}

func createShortURLHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	u := r.Form.Get("Url")
	if u == "" {
        w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("url field cannot be empty"))
		return
	}

	w.Write([]byte("URL field is required"))

	key := shortuuid.New()

	insertMapping(key, u)

	log.Println("url mapped successfully")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("http://localhost:8080/short/%s", key)))
}

func redirectHandler(w http.ResponseWriter, r *http.Request){
	key := chi.URLParam(r, "key")
	if key == "" {
        w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("key field cannot be empty"))
		return
	}

	u := fetchMapping(key)
	if u == "" {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte("url not found"))
    }

	http.Redirect(w,r, u, http.StatusFound)
}
func insertMapping(key string, u string){
	urlMapper.Lock.Lock()
    urlMapper.Mapping[key] = u
    defer urlMapper.Lock.Unlock()
}

func fetchMapping(key string) string{
	urlMapper.Lock.Lock()
    defer urlMapper.Lock.Unlock()
    key, ok := urlMapper.Mapping[key]
    if ok {
       return urlMapper.Mapping[key]
    }

	return ""
}

