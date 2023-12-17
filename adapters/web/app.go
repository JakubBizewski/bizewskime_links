package web

import (
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
	"github.com/gin-gonic/gin"
)

type targetUrlPayload struct {
	TargetUrl string `json:"targetUrl" binding:"required,url"`
}

type WebApp struct {
	shortLinkService *driver.ShortLinkService
	Router           *gin.Engine
}

func CreateWebApp(shortLinkService *driver.ShortLinkService) *WebApp {
	router := gin.Default()

	router.GET("/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		targetUrl, err := shortLinkService.GetTargetUrl(shortCode)
		if err != nil {
			c.String(500, "Something went wrong")
			return
		}

		if targetUrl == "" {
			c.Redirect(302, "/")
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

		shortCode, err := shortLinkService.GenerateShortLink(targetUrlPayload.TargetUrl)
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(200, gin.H{
			"shortCode": shortCode,
		})
	})

	return &WebApp{
		Router: router,
	}
}

func (webApp *WebApp) Run() error {
	return webApp.Router.Run()
}
