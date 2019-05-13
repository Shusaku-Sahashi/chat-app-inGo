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
	temple   *template.Template
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.once.Do(func() {
		h.temple = template.Must(template.ParseFiles(filepath.Join("template", h.filename)))
	})
	h.temple.Execute(w, req)
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
