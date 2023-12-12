package main

import (
	"github.com/JakubBizewski/jakubme_links/adapters/shortLinksDb"
	"github.com/JakubBizewski/jakubme_links/adapters/web"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
)

func main() {
	shortLinkRepository := shortLinksDb.NewMemoryShortLinkRepository()
	shortenedLinksService := driver.NewShortLinkService(shortLinkRepository)

	webApp := web.NewWebApp(shortenedLinksService)

	err := webApp.Run()
	if err != nil {
		panic(err)
	}
}
