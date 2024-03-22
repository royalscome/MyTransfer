package http

import "github.com/gin-gonic/gin"

func (h *Handler) queryOnlineDevices(c *gin.Context) {
	h.svc.QueryOnlineDevices()
}
