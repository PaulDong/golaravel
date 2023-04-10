package golaravel

import (
  "net/http"
  "strconv"
  "github.com/justinas/nosurf"
)

func (g *Golaravel) SessionLoad(next http.Handler) http.Handler {
  g.InfoLog.Println("SessionLoad called")
  return g.Session.LoadAndSave(next)
}

func (g *Golaravel) NoSurf(next http.Handler) http.Handler {
  csrfHandler := nosurf.New(next)
  secure, _ := strconv.ParseBool(g.config.cookie.secure)

  csrfHandler.SetBaseCookie(http.Cookie{
    HttpOnly: true,
    Path: "/",
    Secure: secure,
    SameSite: http.SameSiteStrictMode,
    Domain: g.config.cookie.domain,
  })
  return csrfHandler
}