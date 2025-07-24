package healthcheck

import (
	"den-den-mushi-Go/internal/control/filters"
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/types"
	"den-den-mushi-Go/pkg/util/convert"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) GetApiV1Healthcheck(c *gin.Context, params oapi.GetApiV1HealthcheckParams) {
	h.Log.Debug("GetApiV1Healthcheck called", zap.Any("params", params))

	f := filters.HealthcheckPtySession{
		Hostname:        params.Hostname,
		Ip:              params.Ip,
		Appcode:         params.Appcode,
		Lob:             params.Lob,
		OsType:          params.OsType,
		Status:          params.Status,
		Environment:     params.Environment,
		Country:         params.Country,
		SystemType:      params.SystemType,
		PtySessionState: (*types.PtySessionState)(params.PtySessionState),
		Page:            convert.DerefOr(params.Page, 1),
		PageSize:        convert.DerefOr(params.PageSize, 20),
	}

	results, err := h.Service.getHostsAndAssociatedPtySessions(f)
	if err != nil {
		h.Log.Error("Failed to fetch healthcheck sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, oapi.HealthcheckSessionsResponse{
		HostSessionDetails: results,
	})
}
