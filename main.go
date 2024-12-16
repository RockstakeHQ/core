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

	filebaseAccessKey := os.Getenv("FILEBASE_ACCESS_KEY")
	filebaseSecretKey := os.Getenv("FILEBASE_SECRET_KEY")
	filebaseBucketName := os.Getenv("FILEBASE_BUCKET")

	fixtureStore := db.NewPocketBaseFixtureStore(app)
	fixtureHandler := api.NewFixtureHandler(fixtureStore)

	betStore := db.NewPocketBaseBetStore(app)
	betHandler := api.NewBetHandler(betStore)

	nftStore := db.NewFilebaseStore(filebaseAccessKey, filebaseSecretKey, filebaseBucketName)
	nftHandler := api.NewNFTHandler(nftStore)

	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/api/fixtures", func(e *core.RequestEvent) error {
			return fixtureHandler.HandleGetFixturesByDate(e)
		})

		e.Router.POST("/api/bet", func(e *core.RequestEvent) error {
			return betHandler.HandlePostBet(e)
		})

		e.Router.POST("/api/nft/generate", func(e *core.RequestEvent) error {
			return nftHandler.HandleGenerateNFT(e)
		})
		return e.Next()
	})

	if len(os.Args) == 1 {
		log.Println("Starting server on http://localhost:8090")
		if err := app.Start(); err != nil {
			log.Fatal(err)
		}
		return
	}

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
