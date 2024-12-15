package signals

import (
	"context"
	"github.com/elastic/go-elasticsearch/v8"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"recommendation/server"
	"syscall"
	"time"
)

var shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

type Shutdown struct {
	// pool                  *redis.Pool
	// tracerProvider        *sdktrace.TracerProvider
	serverShutdownTimeout time.Duration
}

func NewShutdown(serverShutdownTimeout time.Duration) (*Shutdown, error) {
	srv := &Shutdown{
		serverShutdownTimeout: serverShutdownTimeout,
	}

	return srv, nil
}

func (s *Shutdown) Shutdown(stopCh <-chan struct{}, svr *server.Server) {
	ctx := context.Background()

	<-stopCh
	ctx, cancel := context.WithTimeout(context.Background(), s.serverShutdownTimeout)
	defer cancel()

	log.Info("Shutting down HTTP/HTTPS server. ", s.serverShutdownTimeout)

	if err := svr.Shutdown(ctx); err != nil {
		log.Warn("HTTP server graceful shutdown failed", err)
	}
	log.Error("Shutdown complete.")
}

var onlyOneSignalHandler = make(chan struct{})

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func SetupSignalHandler() (stopCh <-chan struct{}) {
	close(onlyOneSignalHandler) // panics when called twice

	stop := make(chan struct{})
	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		log.Error("signal received, stopping...")
		close(stop)
		<-c
		log.Error("signal received, exit...")
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func (s *Shutdown) Graceful(stopCh <-chan struct{}, svr *server.Server, client *elasticsearch.Client) {
	ctx := context.Background()

	// wait for the server to gracefully terminate
	<-stopCh
	ctx, cancel := context.WithTimeout(ctx, s.serverShutdownTimeout)
	defer cancel()

	// all calls to /healthz and /readyz will fail from now on
	// atomic.StoreInt32(healthy, 0)
	// atomic.StoreInt32(ready, 0)

	// close cache pool
	// if s.pool != nil {
	// 	_ = s.pool.Close()
	// }

	log.Info("Shutting down HTTP/HTTPS server. ", s.serverShutdownTimeout)

	// There could be a period where a terminating pod may still receive requests. Implementing a brief wait can mitigate this.
	// See: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination
	// the readiness check interval must be lower than the timeout
	// if viper.GetString("level") != "debug" {
	// 	time.Sleep(3 * time.Second)
	// }

	// // stop OpenTelemetry tracer provider
	// if s.tracerProvider != nil {
	// 	if err := s.tracerProvider.Shutdown(ctx); err != nil {
	// 		s.logger.Warn("stopping tracer provider", zap.Error(err))
	// 	}
	// }

	// determine if the GRPC was started
	// if grpcServer != nil {
	// 	s.logger.Info("Shutting down GRPC server", zap.Duration("timeout", s.serverShutdownTimeout))
	// 	grpcServer.GracefulStop()
	// }

	// determine if the http server was started
	if err := svr.Shutdown(ctx); err != nil {
		log.Warn("HTTP server graceful shutdown failed", err)
	}

	log.Error("Shutdown complete.")
}
