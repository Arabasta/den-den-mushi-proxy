package whiteblacklist

import (
	oapi "den-den-mushi-Go/openapi/control"
	"den-den-mushi-Go/pkg/types"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func (h *Handler) GetApiV1WhitelistRegex(c *gin.Context) {
	filters, err := h.Service.GetRegexFilters(types.Whitelist, c)
	if err != nil {
		h.Log.Error("failed to fetch whitelist filters", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	h.Log.Debug("GET filters", zap.Any("filters", filters))
	c.JSON(http.StatusOK, filters)
}

func (h *Handler) PostApiV1WhitelistRegex(c *gin.Context) {
	var req oapi.PostApiV1WhitelistRegexJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter, err := h.Service.CreateRegex(req.Pattern, types.Whitelist, c)
	if err != nil {
		h.Log.Error("failed to create filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusCreated, filter)
}

func (h *Handler) PutApiV1WhitelistRegexId(c *gin.Context, id int) {
	var req oapi.PutApiV1WhitelistRegexIdJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Log.Warn("invalid update body", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	filter, err := h.Service.UpdateRegex(id, req.Pattern, *req.IsEnabled, c)
	if err != nil {
		h.Log.Error("failed to update filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, filter)
}

func (h *Handler) DeleteApiV1WhitelistRegexId(c *gin.Context, id int) {
	filter, err := h.Service.SoftDeleteRegex(id, c)
	if err != nil {
		h.Log.Error("failed to delete filter", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}

	c.JSON(http.StatusOK, filter)
}
