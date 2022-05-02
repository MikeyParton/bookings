package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/MikeyParton/bookings/pkg/config"
	"github.com/MikeyParton/bookings/pkg/handlers"
	"github.com/MikeyParton/bookings/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8082"

var app config.AppConfig

var session *scs.SessionManager

func main() {
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal(err)
	}

	app.IsProduction = false
	app.TemplateCache = tc
	app.UseCache = false

	session = scs.New()
	session.Lifetime = time.Hour * 24
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Persist = true
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	render.NewTemplates(&app)
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	fmt.Printf("starting app on %s", portNumber)
	srv := &http.Server{
		Handler: Routes(&app),
		Addr:    portNumber,
	}

	srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
