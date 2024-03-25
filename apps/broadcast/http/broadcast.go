package http

import (
	"MyTransfer/apps/broadcast"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
)

func (h *Handler) queryOnlineDevices(c *gin.Context) {
	resp := h.svc.QueryOnlineDevices()
	response.Success(c.Writer, resp)
}

func (h *Handler) sendMessageUseUDP(c *gin.Context) {
	ins := broadcast.NewUDPMessage()
	if err := c.Bind(ins); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	fmt.Println(ins)
	if err := h.svc.SendMessageUseUDP(h.c, ins); err != nil {
		response.Failed(c.Writer, err)
		return
	}
	response.Success(c.Writer, "发送成功")
}
