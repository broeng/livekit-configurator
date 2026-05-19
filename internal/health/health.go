package health

import (
	"context"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/broeng/livekit-configurator/internal/livekit"
	"github.com/broeng/livekit-configurator/internal/sleep"

	"github.com/sirupsen/logrus"
)

type HealthController struct {
	logger logrus.FieldLogger
	context context.Context
	client  *livekit.LiveKitClient
	listenPort int
	isReady atomic.Bool
}

func New(logger logrus.FieldLogger, context context.Context, client *livekit.LiveKitClient, listenPort int) *HealthController {
	return &HealthController {
		logger: logger.WithField("component", "health"),
		context: context,
		client: client,
		listenPort: listenPort,
	}
}

func (hc *HealthController) Run() {
	// Register Kubernetes probe endpoints
	http.HandleFunc("/healthz", hc.livenessHandler)
	http.HandleFunc("/readyz", hc.readinessHandler)

	// Prepare Server
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", hc.listenPort),
		Handler: nil,
	}

	// Register shutdown handler
	go func() {
		<-hc.context.Done()
		hc.logger.Info("Shutting down HTTP server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			hc.logger.Errorf("Health HTTP Server shutdown error: %s", err)
		}
	}()

	// Register readiness probe
	go hc.readinessProbe()

	// Start HTTP server
	hc.logger.Infof("Server starting on :%d ...", hc.listenPort)
	if err := server.ListenAndServe(); err != nil {
		if hc.context.Err() == nil {
			hc.logger.Errorf("Failed to set up Health service: %s", err)
		}
	}

	hc.logger.Info("Health HTTP server shut down.")
}

func (hc *HealthController) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (hc* HealthController) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if hc.isReady.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("Not Ready: unable to query LiveKit server"))
	}
}

func (hc* HealthController) readinessProbe() {
	s := sleep.New(hc.context)
	for s.IsRunning() {
		// Probe server endpoint
		hc.isReady.Store(hc.client.TestConnection())
		s.Sleep(30)
	}
}
