package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie || cookie.Value == "" {
		// not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.next.ServeHTTP(w, r)
}

// MustAuth wraps given handler with authHandler
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process
// format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]

	switch action {
	case "login":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to get provider %s: %s", provider, err)
			http.Error(w, msg, http.StatusBadRequest)
		}

		loginURL, err := provider.GetBeginAuthURL(nil, nil)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to GetBeginAuthURL for %s: %s", provider, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		w.Header().Set("Location", loginURL)
		w.WriteHeader(http.StatusTemporaryRedirect)

	case "callback":
		provider, err := gomniauth.Provider(provider)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to get provider %s: %s", provider, err)
			http.Error(w, msg, http.StatusBadRequest)
		}

		creds, err := provider.CompleteAuth(objx.MustFromURLQuery(r.URL.RawQuery))
		if err != nil {
			msg := fmt.Sprintf("Error when trying to complete auth for %s: %s", provider, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		user, err := provider.GetUser(creds)
		if err != nil {
			msg := fmt.Sprintf("Error when trying to get user for %s: %s", provider, err)
			http.Error(w, msg, http.StatusInternalServerError)
		}

		authCookieValue := objx.New(map[string]interface{}{
			"name":       user.Name(),
			"avatar_url": user.AvatarURL(),
			"email":      user.Email(),
			"user_id":    hashUserID(user.Email()),
		}).MustBase64()

		http.SetCookie(w, &http.Cookie{
			Name:  "auth",
			Value: authCookieValue,
			Path:  "/",
		})
		w.Header().Set("Location", "/chat")
		w.WriteHeader(http.StatusTemporaryRedirect)

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}

func hashUserID(email string) string {
	m := md5.New()
	io.WriteString(m, strings.ToLower(email))
	return fmt.Sprintf("%x", m.Sum(nil))
}
