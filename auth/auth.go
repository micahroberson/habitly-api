package auth

import (
  "encoding/json"
  "fmt"
  "net/http"
  "net/url"
  "time"

  "github.com/go-martini/martini"
  "github.com/martini-contrib/sessions"
)

const (
  keyToken     = "auth_token"
  keyEmail     = "auth_email"
)

func AuthProvider() martini.Handler {
  return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
    token = r.Header.Get(keyToken)
    email = r.Header.Get(keyEmail)

    if token != nil && email != nil {
      s.Set(keyToken, token)
      s.Set(keyEmail, token)
    }
  }
}

var LoginRequired martini.Handler = func() martini.Handler {
  return func(s sessions.Session, c martini.Context, w http.ResponseWriter, r *http.Request) {
    token, email := unmarshallTokenAndEmail(s)
    if token == nil || email == nil {
      http.Error(w, "Not Authorized", 401)
    }
  }
}()

func unmarshallTokenAndEmail(s sessions.Session) (t string, e string) {
  if s.Get(keyToken) == nil || s.Get(keyEmail) == nil {
    return
  }

  token := s.Get(keyToken)
  email := s.Get(keyEmail)

  return tokenken, email
}