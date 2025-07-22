package make_change

import (
	"den-den-mushi-Go/internal/control/filters"
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/httpx"
	"den-den-mushi-Go/pkg/middleware"
	"den-den-mushi-Go/pkg/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	Service *Service
	Log     *zap.Logger
}

func (h *Handler) GetApiV1ChangeRequests(c *gin.Context, params oapi.GetApiV1ChangeRequestsParams) {
	h.Log.Debug("GetApiV1ChangeRequests called", zap.Any("params", params))

	authCtx, ok := middleware.GetAuthContext(c.Request.Context())
	if !ok {
		httpx.RespondError(c, http.StatusUnauthorized, "auth context missing", nil, h.Log)
		return
	}

	r := filters.ListCR{
		TicketIDs:         params.TicketIds,
		ImplementorGroups: params.ImplementorGroups,
		LOB:               params.Lob,
		Country:           params.Country,
		StartTime:         params.StartTime,
		EndTime:           params.EndTime,
		PtySessionState:   (*types.PtySessionState)(params.PtySessionState),
		Page:              derefOr(params.Page, 1),
		PageSize:          derefOr(params.PageSize, 20),
	}

	results, err := h.Service.ListChangeRequestsWithSessions(r, authCtx)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "failed to retrieve change requests", err, h.Log)
		return
	}

	c.JSON(http.StatusOK, results)
}

func derefOr[T any](ptr *T, fallback T) T {
	if ptr != nil {
		return *ptr
	}
	return fallback
}
