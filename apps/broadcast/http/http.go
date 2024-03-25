package http

import (
	"MyTransfer/apps"
	"MyTransfer/apps/broadcast"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger/zap"
	"net"
)

var (
	handler = &Handler{}
)

type Handler struct {
	svc broadcast.Service
	c   *net.UDPConn
}

func (h *Handler) Config(c interface{}) {
	conn, ok := c.(*net.UDPConn)
	if !ok {
		zap.L().Errorf("conn error")
		return
	}

	h.svc = apps.GetImpl(broadcast.AppName).(broadcast.Service)
	h.c = conn
}

func (h *Handler) Registry(r gin.IRouter) {
	r.POST("/getOnlineDevices", h.queryOnlineDevices)
	r.POST("/sendMessage", h.sendMessageUseUDP)
}

func (h *Handler) Name() string {
	return broadcast.AppName
}

func init() {
	apps.RegistryGin(handler)
}
