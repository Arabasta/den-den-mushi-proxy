package meta

import (
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Log     *zap.Logger
	Service *Service
}

func (h Handler) GetApiV1Metadata(c *gin.Context) {
	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		h.Log.Error("Auth context missing in request")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "auth context missing"})
		return
	}

	userImplGroups, err := h.Service.getUserImplementorGroups(authCtx)
	if err != nil {
		h.Log.Error("fetch user implementor groups", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user implementor groups"})
		return
	}

	payload := oapi.UserMetadata{
		ImplementorGroups: userImplGroups,
		UserId:            authCtx.UserID,
		OuGroup:           authCtx.OuGroup,
		Roles:             []string{}, // for future use
	}

	if payload.ImplementorGroups == nil {
		payload.ImplementorGroups = []string{}
	}

	c.JSON(http.StatusOK, payload)
}
