package main

import (
	"github.com/JakubBizewski/jakubme_links/adapters/shortLinksDb"
	"github.com/JakubBizewski/jakubme_links/adapters/web"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
)

func main() {
	shortLinkRepository := shortLinksDb.CreateMemoryShortLinkRepository()
	shortLinksService := driver.CreateShortLinkService(shortLinkRepository)

	webApp := web.CreateWebApp(shortLinksService)

	err := webApp.Run()
	if err != nil {
		panic(err)
	}
}
