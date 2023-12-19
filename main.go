package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JakubBizewski/jakubme_links/adapters/sqlite"
	"github.com/JakubBizewski/jakubme_links/adapters/web"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
)

func main() {
	dbPath := os.Getenv("DB_PATH")

	dbSetupErr := sqlite.Setup(dbPath)
	if dbSetupErr != nil {
		panic(dbSetupErr)
	}

	shortLinkRepository := sqlite.CreateShortLinkRepository(dbPath)
	shortLinksService := driver.CreateShortLinkService(shortLinkRepository)

	webApp := web.CreateWebApp(shortLinksService)

	webErr := webApp.Run()
	if webErr != nil {
		panic(webErr)
	}
}
