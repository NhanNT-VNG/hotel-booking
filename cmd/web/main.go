package main

import (
	"log"
	"net/http"
	"time"

	"github.com/NhanNT-VNG/hotel-booking/pkg/config"
	"github.com/NhanNT-VNG/hotel-booking/pkg/handlers"
	"github.com/NhanNT-VNG/hotel-booking/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":3000"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	app.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("Cannot create tempalate cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	repo := handlers.NewRepo(&app)

	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	log.Println("App listen on port", portNumber)

	src := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	err = src.ListenAndServe()
	log.Fatal(err)
}
