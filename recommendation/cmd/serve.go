package cmd

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"net"
	"net/http"
	"path/filepath"
	"recommendation/logger"
	"recommendation/server"
	"recommendation/setting"
	"recommendation/signals"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "search-demo-template",
	Short: "A clean architecture template for Golang Gin services",

	Run: func(cmd *cobra.Command, args []string) {
		var config setting.Configuration

		if err := setting.LoadConfig(&config); err != nil {
			log.Error("loading config file failed.", err)
			return
		}

		gin.SetMode(config.Server.RunMode)

		// init logger
		logger.SetupLogger(
			filepath.Join(config.Log.LogSavePath, config.Log.LogFileName),
			config.Log.MaxSize, config.Log.MaxBackups, config.Log.Compress,
			config.Log.Level)

		// init db
		engine, err := dbInit(&config.Database)
		if err != nil {
			log.Error("init db failed.", err)
			return
		}

		// start http server
		svr := server.NewServer(&config, engine)
		if err := svr.Start(); err != nil {
			log.Error("init server failed.", err)
			return
		}

		// graceful shutdown
		stopCh := signals.SetupSignalHandler()
		sd, _ := signals.NewShutdown(config.App.ServerShutdownTimeout)
		sd.Graceful(stopCh, svr, engine)
	},
}

func Execute() error {
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	return rootCmd.Execute()
}

func dbInit(dbc *setting.DatabaseSettingS) (*elasticsearch.Client, error) {
	// var err error
	config := elasticsearch.Config{
		Addresses: []string{
			dbc.Address,
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   dbc.IdleHost,
			ResponseHeaderTimeout: time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
		},
	}
	return elasticsearch.NewClient(config)
}
