package main

import (
	"fmt"
	"net/http"
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
		panic(fmt.Sprint("AuthorizationErr: %s", err))
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
