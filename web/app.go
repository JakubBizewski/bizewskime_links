package web

import (
	"github.com/JakubBizewski/jakubme_links/shortLinks/application"
	"github.com/JakubBizewski/jakubme_links/shortLinks/domain"
	"github.com/gin-gonic/gin"
)

type targetUrlPayload struct {
	TargetUrl string `json:"targetUrl" binding:"required"`
}

type WebApp struct {
	shortenedLinksService *application.ShortenedLinksService
	router                *gin.Engine
}

func NewWebApp(shortenedLinksService *application.ShortenedLinksService) *WebApp {
	router := gin.Default()

	router.GET("/:shortCode", func(c *gin.Context) {
		shortCodeRaw := c.Param("shortCode")
		shortCode, err := domain.CreateShortCode(shortCodeRaw)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		targetUrl, err := shortenedLinksService.GetLinkByShortCode(shortCode)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		if targetUrl == "" {
			c.JSON(404, gin.H{
				"error": "Not found",
			})
			return
		}

		c.Redirect(302, targetUrl)
	})

	router.POST("/new", func(c *gin.Context) {
		var targetUrlPayload targetUrlPayload
		err := c.BindJSON(&targetUrlPayload)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		targetUrl, err := domain.CreateUrl(targetUrlPayload.TargetUrl)
		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error(),
			})
			return
		}

		link, err := shortenedLinksService.CreateShortenedLink(targetUrl)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"shortCode": link.ShortCode,
		})
	})

	return &WebApp{
		router: router,
	}
}

func (webApp *WebApp) Run() error {
	return webApp.router.Run()
}
