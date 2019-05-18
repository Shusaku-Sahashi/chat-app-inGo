package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
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
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}

	return "", ErrorNoAvatarURL
}
