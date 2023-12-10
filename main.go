package main

import (
	"github.com/JakubBizewski/jakubme_links/shortLinks/application"
	"github.com/JakubBizewski/jakubme_links/shortLinks/infrastructure"
	"github.com/JakubBizewski/jakubme_links/web"
)

func main() {
	domainRepository := infrastructure.NewMemoryDomainRepository()
	shortenedLinksService := application.NewShortenedLinksService(domainRepository)

	webApp := web.NewWebApp(shortenedLinksService)

	err := webApp.Run()
	if err != nil {
		panic(err)
	}
}
