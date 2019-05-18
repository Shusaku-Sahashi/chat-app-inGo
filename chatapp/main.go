package main

import (
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

type templateHandler struct {
	once     sync.Once
	filename string
	temple   *template.Template
}

func (h templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.once.Do(func() {
		h.temple = template.Must(template.ParseFiles(filepath.Join("template", h.filename)))
	})

	data := map[string]interface{}{
		"Host": req.Host,
	}

	name, err := req.Cookie("auth")
	if err == nil {
		data["UserData"] = objx.MustFromBase64(name.Value)
	}
	h.temple.Execute(w, data)
}

func main() {
	gomniauth.SetSecurityKey("1faskfagasgsnnnndiasodff")
	gomniauth.WithProviders(
		google.New("768763966104-fn5svd9ge0mn9v5ie2ntqt0qnivl3ktt.apps.googleusercontent.com", "s90rgodeAOHtDSqgU74aGqFh", "http://localhost:8080/auth/callback/google"),
	)

	r := NewRoom(UserAvatar)
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.HandleFunc("/logout", logoutHandler)
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
