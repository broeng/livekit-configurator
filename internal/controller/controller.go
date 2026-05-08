package controller

import (
	"context"

	"github.com/broeng/livekit-configurator/internal/config"
	"github.com/broeng/livekit-configurator/internal/livekit"
	"github.com/broeng/livekit-configurator/internal/sleep"
	"github.com/broeng/livekit-configurator/internal/types"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger    logrus.FieldLogger
	context   context.Context
	config    config.ConfigT
    client    *livekit.LiveKitClient
}

func New(logger logrus.FieldLogger, context context.Context, config config.ConfigT, client *livekit.LiveKitClient) *Controller {
	c := &Controller{
		logger:    logger.WithField("component", "controller"),
		context:   context,
		config:    config,
		client:    client,
	}
	return c
}

func (c *Controller) isRunning() bool {
	return c.context.Err() == nil
}

func (c *Controller) Reconcile(configDef *types.LiveKitConfigDefinition) {
	s := sleep.New(c.context)
	for c.isRunning() {
		c.logger.Info("Running config reconciliation ...")

		if configDef.SIPConfig != nil {
			c.ReconcileSIPConfig(configDef.SIPConfig)
		}

		if c.config.RunOnce {
			c.logger.Info("Exiting ...")
			break
		} else {
			s.Sleep(c.config.ReconciliationFrequency)
		}
	}
}

func (c *Controller) ReconcileSIPConfig(sipConfig *types.SIPConfig) {
	c.logger.Debug("Reconciling SIP config ...")

	trunkResp, err := c.client.ListSIPTrunks()
	if err != nil {
		c.logger.Warn("failed to list SIP trunks, skipping this round... %s", err)
	} else {
		c.logger.WithField("items", len(trunkResp.Items)).Info("got SIP trunks from server")
	}

}
