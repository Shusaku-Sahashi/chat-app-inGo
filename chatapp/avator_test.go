package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.AvatarURL(client)
	if err != ErrorNoAvatarURL {
		t.Error("値が存在しない場合、エラーになるべきです。")
	}

	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl}
	url, err = authAvatar.AvatarURL(client)
	if err != nil {
		t.Error("エラー無くurlが取得出来るべきです。", err)
	} else if url != testUrl {
		t.Errorf("セットした値と返却された値が異なります。expected: %s, actual: %s", testUrl, url)
	}
}

func TestGravater(t *testing.T) {
	t.Run("対象の画像が取得可能", func(t *testing.T) {
		var gravatar Gravatar
		client := new(client)

		email := "MyEmailAddress@example.com"

		m := md5.New()
		io.WriteString(m, strings.ToLower(email))

		client.userData = map[string]interface{}{"user_id": fmt.Sprintf("%x", m.Sum(nil))}
		url, err := gravatar.AvatarURL(client)

		expected := fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil))

		if err == ErrorNoAvatarURL {
			t.Error("エラー無くurlが取得できる必要があります。", err)
		} else if url != expected {
			t.Error("想定しているURLが取得できるべきです。", "actual: ", url, "expected: ", expected)
		}
	})
}

func TestFileAvatar(t *testing.T) {
	t.Run("テスト用のファイルが取得できる。", func(t *testing.T) {
		var fileSystemAvatar FileSystemAvatar

		filename := filepath.Join("avatars", "abc.jpg")

		ioutil.WriteFile(filename, []byte{}, 0600)
		defer os.Remove(filename)

		client := new(client)

		_, err := fileSystemAvatar.AvatarURL(client)
		if err != ErrorNoAvatarURL {
			t.Error("URLが見つからない場合は、エラーが発生するべきです。")
		}

		client.userData = map[string]interface{}{"user_id": "abc"}

		url, err := fileSystemAvatar.AvatarURL(client)
		if err != nil {
			t.Error("URLが見つからない場合はErrorが発生すべきです。")
		}
		
		if url != "/" + filename {
			t.Errorf("URLが想定と異なる値を返しました。 actual: %s, expected: %s", url, filename)
		}
	})
}
