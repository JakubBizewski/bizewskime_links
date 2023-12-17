package main

import (
	"github.com/JakubBizewski/jakubme_links/adapters/memory"
	"github.com/JakubBizewski/jakubme_links/adapters/web"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
)

func main() {
	shortLinkRepository := memory.CreateMemoryShortLinkRepository()
	shortLinksService := driver.CreateShortLinkService(shortLinkRepository)

	webApp := web.CreateWebApp(shortLinksService)

	err := webApp.Run()
	if err != nil {
		panic(err)
	}
}
