package server

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"recommendation/api"
	"recommendation/logger"
	"recommendation/setting"
)

type Server struct {
	Router *gin.Engine
	Svr    *http.Server
	Config *setting.Configuration
	DB     *elasticsearch.Client
}

func NewServer(cfg *setting.Configuration, client *elasticsearch.Client) *Server {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	srv := &http.Server{
		Addr:           ":" + cfg.Server.HTTPPort,
		Handler:        r,
		ReadTimeout:    cfg.Server.ReadTimeout,
		WriteTimeout:   cfg.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s := &Server{
		Config: cfg,
		Svr:    srv,
		Router: r,
		DB:     client,
	}

	api.SetRouters(r, cfg, client)

	// hooking middleware
	logger.SetLoggerHooking(r)

	// swagger middleware

	return s
}

func (s *Server) Start() error {
	// Timeout: https://adam-p.ca/blog/2022/01/golang-http-server-timeouts/
	go func() {
		log.Info("Starting HTTP Server at :", s.Config.Server.HTTPPort)
		if err := s.Svr.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal("HTTP server expcetpion. ", err)
		}
	}()

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	// TODO: add code
	return s.Svr.Shutdown(ctx)
}
