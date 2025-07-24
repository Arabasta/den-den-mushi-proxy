package make_change

import (
	"den-den-mushi-Go/internal/control/filters"
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/httpx"
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

func (h *Handler) GetApiV1ChangeRequests(c *gin.Context, params oapi.GetApiV1ChangeRequestsParams) {
	h.Log.Debug("GetApiV1ChangeRequests called", zap.Any("params", params))

	r := filters.ListCR{
		TicketIDs:         params.TicketIds,
		ImplementorGroups: params.ImplementorGroups,
		LOB:               params.Lob,
		Country:           params.Country,
		StartTime:         params.StartTime,
		EndTime:           params.EndTime,
		PtySessionState:   (*types.PtySessionState)(params.PtySessionState),
		Page:              convert.DerefOr(params.Page, 1),
		PageSize:          convert.DerefOr(params.PageSize, 20),
	}

	results, err := h.Service.ListChangeRequestsWithSessions(r, c)
	if err != nil {
		httpx.RespondError(c, http.StatusInternalServerError, "failed to retrieve change requests", err, h.Log)
		return
	}

	c.JSON(http.StatusOK, results)
}
