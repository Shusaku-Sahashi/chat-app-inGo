package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
)

type httpHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.once.Do(func() {
		h.templ = template.Must(template.ParseFiles(filepath.Join("template", h.filename)))
	})
	h.templ.Execute(w, req)
}

func main() {
	r := NewRoom()
	http.Handle("/", MustAuth(&httpHandler{filename: "chat.html"}))
	http.Handle("/room", r)
	go r.Run()

	//TODO: export conf to JSON or YAML file
	server := http.Server{

		Addr: "127.0.0.1:8080",

		Handler: nil,
	}

	log.Println("Starting web server on", server.Addr)

	if err := server.ListenAndServe(); err != nil {

		log.Fatal("listenServerError:", err)

	}
}
