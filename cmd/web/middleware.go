package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

func NoSurf(next http.Handler) http.Handler {
	csrfHandle := nosurf.New(next)
	csrfHandle.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandle
}

func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}
