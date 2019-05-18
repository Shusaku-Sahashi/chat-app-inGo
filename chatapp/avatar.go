package main

import "errors"

var ErrorNoAvatarURL = errors.New("chat: アバターのURLを返すことが出来ない。")

type Avatar interface {
	AvatarURL(*client) (string, error)
}

type AuthAvatar struct{}

var UserAvatar AuthAvatar

func (AuthAvatar) AvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrorNoAvatarURL
}
