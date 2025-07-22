package httpx

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RespondError(c *gin.Context, status int, msg string, err error, log *zap.Logger) {
	if err != nil {
		log.Error(msg, zap.Error(err))
	}
	c.AbortWithStatusJSON(status, gin.H{"error": msg})
}
