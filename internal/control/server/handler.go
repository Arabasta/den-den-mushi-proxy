package server

import (
	"den-den-mushi-Go/internal/control/app/healthcheck"
	"den-den-mushi-Go/internal/control/app/iexpress"
	"den-den-mushi-Go/internal/control/app/make_change"
	"den-den-mushi-Go/internal/control/app/meta"
	"den-den-mushi-Go/internal/control/app/pty_token"
	"den-den-mushi-Go/internal/control/app/whiteblacklist"
	oapi "den-den-mushi-Go/openapi/control"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	PtyHandler            *pty_token.Handler
	MakeChangeHandler     *make_change.Handler
	WhiteBlacklistHandler *whiteblacklist.Handler
	HealthcheckHandler    *healthcheck.Handler
	IExpressHandler       *iexpress.Handler
	MetaHandler           *meta.Handler
}

func (h *MasterHandler) GetApiV1Iexpress(c *gin.Context, params oapi.GetApiV1IexpressParams) {
	h.IExpressHandler.GetApiV1IExpress(c, params)
}

func (h *MasterHandler) GetApiV1IexpressRequestId(c *gin.Context, requestId string) {
	h.IExpressHandler.GetApiV1IexpressRequestId(c, requestId)
}

// Forwarding methods (required by oapi.ServerInterface)

func (h *MasterHandler) PostApiV1PtyTokenStart(c *gin.Context) {
	h.PtyHandler.PostApiV1PtyTokenStart(c)
}

func (h *MasterHandler) PostApiV1PtyTokenJoin(c *gin.Context) {
	h.PtyHandler.PostApiV1PtyTokenJoin(c)
}

func (h *MasterHandler) GetApiV1ChangeRequests(c *gin.Context, params oapi.GetApiV1ChangeRequestsParams) {
	h.MakeChangeHandler.GetApiV1ChangeRequests(c, params)
}

func (h *MasterHandler) GetApiV1WhitelistRegex(c *gin.Context) {
	h.WhiteBlacklistHandler.GetApiV1WhitelistRegex(c)
}

func (h *MasterHandler) PostApiV1WhitelistRegex(c *gin.Context) {
	h.WhiteBlacklistHandler.PostApiV1WhitelistRegex(c)
}

func (h *MasterHandler) PutApiV1WhitelistRegexId(c *gin.Context, id int) {
	h.WhiteBlacklistHandler.PutApiV1WhitelistRegexId(c, id)
}

func (h *MasterHandler) DeleteApiV1WhitelistRegexId(c *gin.Context, id int) {
	h.WhiteBlacklistHandler.DeleteApiV1WhitelistRegexId(c, id)
}

func (h *MasterHandler) GetApiV1BlacklistRegex(c *gin.Context) {
	h.WhiteBlacklistHandler.GetApiV1BlacklistRegex(c)
}

func (h *MasterHandler) PostApiV1BlacklistRegex(c *gin.Context) {
	h.WhiteBlacklistHandler.PostApiV1BlacklistRegex(c)
}

func (h *MasterHandler) PutApiV1BlacklistRegexId(c *gin.Context, id int) {
	h.WhiteBlacklistHandler.PutApiV1BlacklistRegexId(c, id)
}

func (h *MasterHandler) DeleteApiV1BlacklistRegexId(c *gin.Context, id int) {
	h.WhiteBlacklistHandler.DeleteApiV1BlacklistRegexId(c, id)
}

func (h *MasterHandler) GetApiV1Healthcheck(c *gin.Context, params oapi.GetApiV1HealthcheckParams) {
	h.HealthcheckHandler.GetApiV1Healthcheck(c, params)
}

func (h *MasterHandler) GetApiV1Metadata(c *gin.Context) {
	h.MetaHandler.GetApiV1Metadata(c)
}
