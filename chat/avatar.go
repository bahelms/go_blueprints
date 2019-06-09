package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

// ErrNoAvatarURL is the error that is returned when the Avatar instance
// is unable to provide a URL
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar represents types capable of representing user profile pictures
type Avatar interface {
	// GetAvatarURL gets the avatar URL for the specified client, or returns
	// an error if something goes wrong.
	//
	// ErrNoAvatarURL is returned if the object is unable to get a URL for
	// the specified client.
	GetAvatarURL(c *client) (string, error)
}

// AuthAvatar stuff
type AuthAvatar struct{}

// UseAuthAvatar stuff
var UseAuthAvatar AuthAvatar

// GetAvatarURL stuff
func (AuthAvatar) GetAvatarURL(c *client) (string, error) {
	url, ok := c.userData["avatar_url"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	urlStr, ok := url.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	return urlStr, nil
}

// GravatarAvatar stuff
type GravatarAvatar struct{}

// UseGravatar implementation
var UseGravatar GravatarAvatar

const gravatarBaseURL = "//www.gravatar.com/avatar"

// GetAvatarURL stuff
func (GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	email, ok := c.userData["email"]
	if !ok {
		return "", ErrNoAvatarURL
	}
	emailStr, ok := email.(string)
	if !ok {
		return "", ErrNoAvatarURL
	}
	m := md5.New()
	io.WriteString(m, strings.ToLower(emailStr))
	return fmt.Sprintf(gravatarBaseURL+"/%x", m.Sum(nil)), nil
}
