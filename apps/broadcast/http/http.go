package http

import (
	"MyTransfer/apps"
	"MyTransfer/apps/broadcast"
	"github.com/gin-gonic/gin"
)

var (
	handler = &Handler{}
)

type Handler struct {
	svc broadcast.Service
}

func (h *Handler) Config() {
	h.svc = apps.GetImpl(broadcast.AppName).(broadcast.Service)
}

func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/getOnlineDevices", h.queryOnlineDevices)
}

func (h *Handler) Name() string {
	return broadcast.AppName
}

func init() {
	apps.RegistryGin(handler)
}
