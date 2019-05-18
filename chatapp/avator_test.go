package main

import "testing"

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
