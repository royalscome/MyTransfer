package protocol

import (
	"MyTransfer/apps"
	"MyTransfer/conf"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/infraboard/mcube/logger"
	"github.com/infraboard/mcube/logger/zap"
	"net/http"
	"time"
)

func NewHttpService() *HttpService {
	r := gin.Default()

	server := &http.Server{
		ReadHeaderTimeout: 60 * time.Second,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    1 << 20, //1M
		Addr:              conf.C().HTTP.HttpAddr(),
		Handler:           r,
	}
	return &HttpService{
		server: server,
		l:      zap.L().Named("HTTP Service"),
		r:      r,
	}
}

type HttpService struct {
	server *http.Server
	l      logger.Logger
	r      gin.IRouter
}

func (s *HttpService) Start() error {
	// 加载Handler
	apps.InitGin(s.r)

	// 已加载App的日志信息
	apps := apps.LoadedGinApps()
	s.l.Infof("loaded gin apps: %v", apps)

	// 该操作是阻塞的,监听端口，等待请求
	// 如果是服务的正常关闭
	if err := s.server.ListenAndServe(); err != nil {
		//if errors.Is(err, http.ErrServerClosed) {
		//	s.l.Info("service stopped success")
		//}
		if err == http.ErrServerClosed {
			s.l.Info("service stopped success")
			return nil
		}
		return fmt.Errorf("start service error, %s", err.Error())
	}
	return nil
}

func (s *HttpService) Stop() {
	s.l.Info("start graceful shutdown")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.l.Warnf("shut down http service error, %s", err)
	}
}
