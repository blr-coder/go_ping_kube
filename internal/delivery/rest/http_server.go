package rest

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type PingServer struct {
	server      *http.Server
	pingHandler *PingHandler
}

func NewPingServer(handler *PingHandler) *PingServer {
	return &PingServer{
		pingHandler: handler,
	}
}

func (s *PingServer) Start(port string) error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      s.handleHttp(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	logrus.Infof("SERVER STARTED ON PORT: %s", port)

	return s.server.ListenAndServe()
}

func (s *PingServer) handleHttp() http.Handler {
	handler := http.NewServeMux()

	handler.HandleFunc("/api/v1/ping", s.pingHandler.Ping)
	handler.HandleFunc("/api/v1/get", s.pingHandler.Get)
	handler.HandleFunc("/api/v1/all", s.pingHandler.All)

	return handler
}
