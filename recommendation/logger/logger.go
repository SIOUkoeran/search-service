package logger

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var logger *lumberjack.Logger

func SetupLogger(filename string, maxsize int, maxbackup int, compress bool, level string, client *elasticsearch.Client) {
	lvl, err := log.ParseLevel(level)
	if err != nil {
		log.SetLevel(lvl)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	// use lumberjack to write to implement rotation.
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

	// use elasticsearch hooking
	log.AddHook(NewElasticsearchHook(client, "log"))
}

func Close() {
	logger.Close()
}

func WithTrace(ctx *gin.Context) *log.Entry {
	fields := log.Fields{}
	if len(ctx.GetString("X-Trace-ID")) > 0 {
		fields["trace_id"] = ctx.GetString("X-Trace-ID")
	}
	if len(ctx.GetString("X-Span-ID")) > 0 {
		fields["span_id"] = ctx.GetString("X-Span-ID")
	}
	return log.WithFields(fields)
}

func SetLoggerHooking(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()

		log.WithFields(log.Fields{
			"method":  c.Request.Method,
			"path":    c.Request.URL.Path,
			"status":  status,
			"latency": latency,
		}).Info("request completed")

		// Add trace and span IDs if available
		if traceEntry := WithTrace(c); traceEntry != nil {
			for k, v := range traceEntry.Data {
				traceEntry = traceEntry.WithField(k, v)
			}
		}

	})
}
