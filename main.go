package main

import (
	"os"

	_ "github.com/mattn/go-sqlite3"

	"github.com/JakubBizewski/jakubme_links/adapters/encryption"
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

	encryptionService := encryption.CreateAESEncryptionService(os.Getenv("ENCRYPTION_KEY"))
	shortLinkRepository := sqlite.CreateShortLinkRepository(dbPath)
	shortLinksService := driver.CreateShortLinkService(shortLinkRepository)

	webApp := web.CreateWebApp(shortLinksService, encryptionService)

	webErr := webApp.Run()
	if webErr != nil {
		panic(webErr)
	}
}
