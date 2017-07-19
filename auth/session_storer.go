package auth


import (
	"net/http"

	"github.com/gorilla/sessions"
	"gopkg.in/authboss.v1"
	"github.com/toasterson/mozaik/logger"
)

const sessionCookieName = "mozaik"

var sessionStore *sessions.CookieStore

type SessionStorer struct {
	w http.ResponseWriter
	r *http.Request
}

func NewSessionStorer(w http.ResponseWriter, r *http.Request) authboss.ClientStorer {
	return &SessionStorer{w, r}
}

func SetSessionStore(store *sessions.CookieStore) {
	sessionStore = store
}

func GetSessionStore() *sessions.CookieStore{
	return sessionStore
}

func (s SessionStorer) Get(key string) (string, bool) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		logger.Error(err)
		return "", false
	}

	strInf, ok := session.Values[key]
	if !ok {
		return "", false
	}

	str, ok := strInf.(string)
	if !ok {
		return "", false
	}

	return str, true
}

func (s SessionStorer) Put(key, value string) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		logger.Error(err)
		return
	}

	session.Values[key] = value
	session.Save(s.r, s.w)
}

func (s SessionStorer) Del(key string) {
	session, err := sessionStore.Get(s.r, sessionCookieName)
	if err != nil {
		logger.Error(err)
		return
	}

	delete(session.Values, key)
	session.Save(s.r, s.w)
}