package pty_token

import (
	"den-den-mushi-Go/internal/control/dto"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RegisterRoutes(r *gin.RouterGroup, s *Service, log *zap.Logger) {
	issr := r.Group("/api/v1/pty_token")
	issr.POST("/start", startHandler(s, log))
	issr.POST("/join", joinHandler(s, log))
}

func joinHandler(s *Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.JoinRequest
		log.Debug("Join request received", zap.Any("request", body))
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error("Failed to bind JSON", zap.Error(err))
			return
		}

		// load keycloak ctx
		body.UserId = c.Keys["user_id"].(string)
		body.OuGroups = c.Keys["ou_groups"].([]string)

		t, p, err := s.mintJoinToken(&body)
		if err != nil {
			log.Error("Failed to mint join token", zap.Error(err))
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Join token minted successfully", zap.String("jwt", t), zap.String("userID", body.UserId))

		c.JSON(200, gin.H{"token": t, "proxyUrl": p})

	}
}

func startHandler(s *Service, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var body dto.StartRequest
		log.Debug("Start request received", zap.Any("request", body))
		err := c.ShouldBindJSON(&body)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			log.Error("Failed to bind JSON", zap.Error(err))
			return
		}

		// load keycloak ctx
		body.UserId = c.Keys["user_id"].(string)
		body.OuGroups = c.Keys["ou_groups"].([]string)

		t, p, err := s.mintStartToken(&body)
		if err != nil {
			log.Error("Failed to mint start token", zap.Error(err))
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		log.Debug("Start token minted successfully", zap.String("jwt", t), zap.String("userID", body.UserId))

		c.JSON(200, gin.H{"token": t, "proxyUrl": p})
	}
}
