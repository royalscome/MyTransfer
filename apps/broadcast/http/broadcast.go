package http

import (
	"MyTransfer/apps/broadcast"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/http/response"
)

type Option interface {
	Apply(*response.Data)
}

func newFuncOption(f func(*response.Data)) Option {
	return &optionFunc{
		f: f,
	}
}

type optionFunc struct {
	f func(data *response.Data)
}

func (s *optionFunc) Apply(resp *response.Data) {
	s.f(resp)
}

func (h *Handler) queryOnlineDevices(c *gin.Context) {
	resp := h.svc.QueryOnlineDevices()
	response.Success(c.Writer, resp, withResponseCode(200))
}

func (h *Handler) sendMessageUseUDP(c *gin.Context) {
	ins := broadcast.NewUDPMessage()
	if err := c.Bind(ins); err != nil {
		response.Failed(c.Writer, err, withResponseCode(500))
		return
	}
	if err := h.svc.SendMessageUseUDP(h.c, ins); err != nil {
		response.Failed(c.Writer, err, withResponseCode(500))
		return
	}
	response.Success(c.Writer, "发送成功", withResponseCode(200))
}

func withResponseCode(code int) Option {
	return newFuncOption(func(data *response.Data) {
		*data.Code = code
	})
}
