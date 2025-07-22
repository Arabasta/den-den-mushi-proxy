package server

import (
	"den-den-mushi-Go/internal/control/make_change"
	"den-den-mushi-Go/internal/control/pty_token"
	"den-den-mushi-Go/internal/control/whiteblacklist"
	oapi "den-den-mushi-Go/openapi/control"
	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	PtyHandler            *pty_token.Handler
	MakeChangeHandler     *make_change.Handler
	WhiteBlacklistHandler *whiteblacklist.Handler
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
