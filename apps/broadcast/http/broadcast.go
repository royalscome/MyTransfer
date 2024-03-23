package http

import (
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
)

func (h *Handler) queryOnlineDevices(c *gin.Context) {
	resp := h.svc.QueryOnlineDevices()
	response.Success(c.Writer, resp)
}
