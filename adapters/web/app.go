package web

import (
	"net/http"

	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
	"github.com/gin-gonic/gin"
)

type targetURLPayload struct {
	TargetURL string `json:"targetURL" binding:"required,url"`
}

type App struct {
	Router *gin.Engine
}

func CreateWebApp(shortLinkService *driver.ShortLinkService) *App {
	router := gin.Default()

	router.GET("/:shortCode", func(c *gin.Context) {
		shortCode := c.Param("shortCode")
		targetURL, err := shortLinkService.GetTargetURL(shortCode)
		if err != nil {
			c.String(http.StatusInternalServerError, "Something went wrong")
			return
		}

		if targetURL == "" {
			c.Redirect(http.StatusFound, "/")
			return
		}

		c.Redirect(http.StatusFound, targetURL)
	})

	router.POST("/new", func(c *gin.Context) {
		var payload targetURLPayload
		err := c.BindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		shortCode, err := shortLinkService.GenerateShortLink(payload.TargetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"shortCode": shortCode,
		})
	})

	return &App{
		Router: router,
	}
}

func (webApp *App) Run() error {
	return webApp.Router.Run()
}
