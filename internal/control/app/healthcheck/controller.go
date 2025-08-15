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

	results, totalCount, err := h.Service.getHostsAndAssociatedPtySessions(f, c)
	if err != nil {
		h.Log.Error("Failed to fetch healthcheck sessions", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//lol
	var response struct {
		TotalCount *int64                           `json:"total_count,omitempty"`
		Items      oapi.HealthcheckSessionsResponse `json:"items"`
	}

	response.Items = oapi.HealthcheckSessionsResponse{
		HostSessionDetails: results,
	}
	response.TotalCount = &totalCount

	if len(*results) == 0 {
		h.Log.Debug("No hosts found for the given filter", zap.Any("filter", f))
		c.JSON(http.StatusOK, response)
		return
	}

	c.JSON(http.StatusOK, response)
}
