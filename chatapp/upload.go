package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/stretchr/objx"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	authCookieValue, err := r.Cookie("auth")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	file, header, err := r.FormFile("avatarFile")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	data, err := ioutil.ReadAll(file)

	defer file.Close()

	if err != nil {
		fmt.Fprint(w, err)
		return
	}

	if userID, ok := objx.MustFromBase64(authCookieValue.Value)["user_id"]; ok {
		if userIDstr, ok := userID.(string); ok {
			filename := filepath.Join("avatars", fmt.Sprintf("%s.%s", userIDstr, header.Filename))
			if err := ioutil.WriteFile(filename, data, 0600); err != nil {
				fmt.Errorf("書き込みに失敗しました。err: %s", err)
				fmt.Fprint(w, "失敗しました。")
			}
		}
	}
	fmt.Fprint(w, "成功")
}
