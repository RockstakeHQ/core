package main

import (
	"log"
	"os"
	"rockstake-core/api"
	"rockstake-core/db"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

func main() {
	app := pocketbase.New()

	// Inițializăm store-ul și handler-ul
	fixtureStore := db.NewPocketBaseFixtureStore(app)
	fixtureHandler := api.NewFixtureHandler(fixtureStore)

	// Setăm rutele
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		// Adăugăm ruta pentru fixtures
		e.Router.GET("/api/fixtures", func(e *core.RequestEvent) error {
			return fixtureHandler.HandleGetFixturesByDate(e)
		})
		return e.Next()
	})

	// Dacă nu sunt argumente, pornim direct serverul
	if len(os.Args) == 1 {
		log.Println("Starting server on http://localhost:8090")
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
		return
	}

	// Altfel, lăsăm PocketBase să gestioneze comenzile
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
