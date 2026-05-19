package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/broeng/livekit-configurator/internal/config"
	"github.com/broeng/livekit-configurator/internal/controller"
	"github.com/broeng/livekit-configurator/internal/health"
	client "github.com/broeng/livekit-configurator/internal/livekit"
	"github.com/broeng/livekit-configurator/internal/types"

	"github.com/sirupsen/logrus"
)

var version = "unreleased"

func main() {
	logger := logrus.New()

	config, err := config.LoadConfig("LIVEKIT_CONF_")
	if err != nil {
		logger.Fatalf("could not parse options: %s", err)
		os.Exit(1)
	}

	if config.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if level, err := logrus.ParseLevel(config.LogLevel); err != nil {
		logger.Fatalf("could not set log level to %s: %s", config.LogLevel, err)
	} else {
		logger.SetLevel(level)
	}

	logger.WithField("version", version).Info("starting LiveKit Configurator")

	// Prepare a default context
	ctx, stopApp := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopApp()

	// Attempt to parse the provided config
	configDef, err := types.ParseConfigDefinition(config.ConfigPath)
	if err != nil {
		logger.Fatalf("could not parse config definition: %s", err)
		os.Exit(1)
	}
	logger.Info("Parsed LiveKit config definition")

	// Prepare client for livekit server
	livekitClient := client.New(logger, ctx, config)

	// Prepare controller handling config reconciliation
	controller := controller.New(logger, ctx, config, livekitClient)

	// Prepare Health Service
	healthController := health.New(logger, ctx, livekitClient, config.HealthListenPort)

	// Prepare exit code
	exitCode := 0

	// Run controllers and await clean exits
	var wg sync.WaitGroup
	wg.Go(healthController.Run)
	wg.Go(func() {
		if controller.Reconcile(configDef) == false {
			exitCode = 1
		}
		stopApp()
	})
	wg.Wait()

	logger.Info("Exiting application ...")
	os.Exit(exitCode)
}
