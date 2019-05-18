package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if _, err := req.Cookie("auth"); err == http.ErrNoCookie {
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

// TODO: 認証はGothに変更する。
func loginHandler(w http.ResponseWriter, req *http.Request) {
	seg := strings.Split(req.URL.Path, "/")
	action := seg[2]
	provider := seg[3]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました。", provider, err)
		}
		url, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			log.Fatalln("GetBeginAuthURLの呼び出し中にエラーが発生しました。", provider, err )
		}
		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusTemporaryRedirect)
	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			log.Fatalln("認証プロバイダーの取得に失敗しました。", provider, err)
		}
		token, err := provider.CompleteAuth(objx.MustFromURLQuery(req.URL.RawQuery))
		if err != nil {
			log.Fatalln("認証を完了出来ませんでした。", provider, err)
		}
		user, err := provider.GetUser(token)
		if err != nil {
			log.Fatalln("ユーザの取得にしっぱいしました。", provider, err)
		}
		// cookieをBase64変換した文字列を登録
		authCookieValue := objx.New(map[string]interface{}{
			"name": user.Name(),
			"avatar_url": user.AvatarURL(),
			"email": user.Email(),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name: "auth",
			Value: authCookieValue,
			Path: "/",
		})

		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です。", action)
	}
}

func logoutHandler(w http.ResponseWriter, req *http.Request) {
	authCookie, err := req.Cookie("auth");
	if err != nil { return }

	authCookie.MaxAge = 0

	http.SetCookie(w, authCookie)

	w.Header().Set("Location", "/login")
	w.WriteHeader(http.StatusTemporaryRedirect)
}