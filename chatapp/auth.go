package main

import (
	"fmt"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie("access_token"); err == http.ErrNoCookie {
		// TODO: check authentication
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(fmt.Sprintf("AuthorizationErr: %s", err))
	} else {
		// Pass Authentication
		h.next.ServeHTTP(w, req)
	}
}

/*
MustAuth is decoretor for handler that needed to authentication.
*/
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{
		next: handler,
	}
}

func loginHandler(w http.ResponseWriter, req *http.Request) {
	seg := strings.Split(req.URL.Path, "/")
	action := seg[2]
	provider := seg[3]

	switch action {
	case "login":
		fmt.Fprint(w, "TODO:Login処理", provider)
	case "callback":
		fmt.Fprint(w, "TODO:Callback処理", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です。", action)
	}
}
