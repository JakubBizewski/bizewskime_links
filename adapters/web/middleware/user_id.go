package middleware

import (
	"github.com/JakubBizewski/jakubme_links/domain/ports/driven"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UserID(encryptionService driven.EncryptionService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := getUserID(c, encryptionService)
		userIDEncrypted, err := encryptionService.Encrypt(userID)
		if err != nil {
			return
		}

		c.SetCookie("user-id", userIDEncrypted, 0, "/", "", false, true)
		c.Set("user_id", userID)

		c.Next()
	}
}

func getUserID(c *gin.Context, encryptionService driven.EncryptionService) string {
	cookieEncryptedUserID, err := c.Cookie("user-id")
	if err != nil {
		return uuid.NewString()
	}

	userID, err := encryptionService.Decrypt(cookieEncryptedUserID)
	if err != nil {
		return uuid.NewString()
	}

	return userID
}
