package web

import (
	"errors"
	"net/http"

	"github.com/JakubBizewski/jakubme_links/adapters/web/middleware"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
	"github.com/JakubBizewski/jakubme_links/domain/ports/driver"
	"github.com/gin-gonic/gin"
)

type targetURLPayload struct {
	TargetURL string `json:"targetURL" binding:"required,url"`
}

type App struct {
	Router *gin.Engine
}

func CreateWebApp(shortLinkService *driver.ShortLinkService, encryptionService driven.EncryptionService) *App {
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

	router.POST("/new", middleware.UserID(encryptionService), func(c *gin.Context) {
		var payload targetURLPayload
		err := c.BindJSON(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		shortCode, err := shortLinkService.GenerateShortLink(payload.TargetURL)
		if errors.Is(err, driver.ErrUniqueShortCodeGenerationFailed) {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate unique short code",
			})
			return
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Something went wrong",
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
