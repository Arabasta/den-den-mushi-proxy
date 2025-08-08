package iexpress

import (
	"den-den-mushi-Go/internal/control/filters"
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/httpx"
	"den-den-mushi-Go/pkg/util/convert"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) GetApiV1IExpress(c *gin.Context, params oapi.GetApiV1IexpressParams) {
	h.Log.Debug("GetApiV1IExpress called", zap.Any("params", params))

	f := filters.ListIexpress{
		TicketIDs:         params.RequestIds,
		Requestor:         params.Requestor,
		ImplementorGroups: nil,
		LOB:               params.Lob,
		Country:           params.OriginCountry,
		AppImpacted:       params.AppImpacted,
		StartTime:         params.StartTime,
		EndTime:           params.EndTime,
		Page:              convert.DerefOr(params.Page, 1),
		PageSize:          convert.DerefOr(params.PageSize, 20),
		IsGetTotalCount:   convert.DerefOr(params.TotalCount, false),
	}

	res, err := h.Service.ListExpressRequests(f, c)
	if err != nil {
		h.Log.Error("Failed to fetch IExpress records", zap.Error(err))
		httpx.RespondError(c, http.StatusInternalServerError, "failed to retrieve iexpress requests", err, h.Log)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetApiV1IexpressRequestId(c *gin.Context, id string) {
	h.Log.Debug("GetApiV1IexpressRequestId called", zap.String("id", id))

	res, err := h.Service.ListExpressRequestDetails(id, c)
	if err != nil {
		h.Log.Error("Failed to fetch IExpress details", zap.Error(err))
		httpx.RespondError(c, http.StatusInternalServerError, "failed to retrieve iexpress requests", err, h.Log)
		return
	}

	c.JSON(http.StatusOK, res)
}
