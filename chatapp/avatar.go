package main

import (
	"errors"
	"fmt"
)

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

type Gravatar struct{}

var UserGravatar Gravatar

func (Gravatar) AvatarURL(c *client) (string, error) {
	if userId, ok := c.userData["user_id"]; ok {
		if userIdStr, ok := userId.(string); ok {
			return fmt.Sprintf("//www.gravatar.com/avatar/%s", userIdStr), nil
		}
	}

	return "", ErrorNoAvatarURL
}
