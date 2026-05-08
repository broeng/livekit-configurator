package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/broeng/livekit-configurator/internal/config"
	"github.com/broeng/livekit-configurator/internal/controller"
	client "github.com/broeng/livekit-configurator/internal/livekit"
	"github.com/broeng/livekit-configurator/internal/types"

	"github.com/sirupsen/logrus"
	"github.com/stevenroose/gonfig"

)

var version = "unreleased"

func main() {
	logger := logrus.New()

	if err := gonfig.Load(&config.Config, gonfig.Conf{
		EnvPrefix:         "LIVEKIT_CONFIGURATOR_",
		FlagIgnoreUnknown: false,
	}); err != nil {
		logger.Fatalf("could not parse options: %s", err)
		os.Exit(1)
	}

	if config.Config.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if level, err := logrus.ParseLevel(config.Config.LogLevel); err != nil {
		logger.Fatalf("could not set log level to %s: %s", config.Config.LogLevel, err)
	} else {
		logger.SetLevel(level)
	}

	logger.WithField("version", version).Info("starting LiveKit Configurator")

	// Prepare a default context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Attempt to parse the provided config
	configDef, err := types.ParseConfigDefinition(config.Config.ConfigPath)
	if err != nil {
		logger.Fatalf("could not parse config definition: %s", err)
		os.Exit(1)
	}
	logger.Info("Parsed LiveKit config definition")

	// Prepare client for livekit server
	livekitClient := client.New(logger, ctx, config.Config)

	// Prepare controller handling config reconciliation
	controller := controller.New(logger, ctx, config.Config, livekitClient)

	controller.Reconcile(configDef)

}
