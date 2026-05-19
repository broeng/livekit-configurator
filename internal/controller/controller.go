package controller

import (
	"context"

	"github.com/broeng/livekit-configurator/internal/config"
	"github.com/broeng/livekit-configurator/internal/conversions"
	"github.com/broeng/livekit-configurator/internal/livekit"
	"github.com/broeng/livekit-configurator/internal/sleep"
	"github.com/broeng/livekit-configurator/internal/types"
	"github.com/sirupsen/logrus"
)

type Controller struct {
	logger  logrus.FieldLogger
	context context.Context
	config  *config.Config
	client  *livekit.LiveKitClient
}

func New(logger logrus.FieldLogger, context context.Context, config *config.Config, client *livekit.LiveKitClient) *Controller {
	c := &Controller{
		logger:  logger.WithField("component", "controller"),
		context: context,
		config:  config,
		client:  client,
	}
	return c
}

func (c *Controller) isRunning() bool {
	return c.context.Err() == nil
}

func (c *Controller) shouldOverwrite() bool {
	return *c.config.ReconciliationMode == config.Overwrite
}
func (c *Controller) shouldMerge() bool {
	return *c.config.ReconciliationMode == config.Merge
}
func (c *Controller) shouldAssert() bool {
	return *c.config.ReconciliationMode == config.Assert
}

func (c *Controller) Reconcile(configDef *types.LiveKitConfigDefinition) bool {
	s := sleep.New(c.context)
	for c.isRunning() {
		c.logger.Info("Running config reconciliation ...")

		if configDef.SIPConfig != nil {
			result := c.ReconcileSIPConfig(configDef.SIPConfig)
			if c.shouldAssert() && result == false {
				return false
			}
		}

		if c.config.RunOnce {
			c.logger.Info("Exiting ...")
			break
		} else {
			s.Sleep(c.config.ReconciliationFrequency)
		}
	}
	c.logger.Info("Reconciliation shut down")
	return true
}

func (c *Controller) ReconcileSIPConfig(sipConfig *types.SIPConfig) bool {
	c.logger.Debug("Reconciling SIP config ...")

	if sipConfig.Trunks != nil {
		trunkConfig := sipConfig.Trunks
		if c.reconcileSIPInboundTrunks(trunkConfig.InboundTrunks) == false && c.shouldAssert() {
			return false
		}
		if c.reconcileSIPOutboundTrunks(trunkConfig.OutboundTrunks) == false && c.shouldAssert() {
			return false
		}
	}
	if sipConfig.DispatchRules != nil {
		if c.reconcileSIPDispatchRules(sipConfig.DispatchRules) == false && c.shouldAssert() {
			return false
		}
	}

	return true
}

func (c *Controller) reconcileSIPInboundTrunks(trunks []*types.SIPInboundBridgeRequest) bool {
	if len(trunks) == 0 {
		return true
	}
	c.logger.Debug("Reconciling SIP inbound trunks ...")
	resp, err := c.client.ListSIPInboundTrunk()
	if err != nil {
		c.logger.Warnf("failed to list SIP inbound trunks, skipping this round... %s", err)
		return false
	}
	c.logger.WithField("items", len(resp.Items)).Info("got SIP inbound trunks from server")

	foundInbounds := make([]bool, len(trunks), len(trunks))
	for trunkIdx := range resp.Items {
		trunk := resp.Items[trunkIdx]
		c.logger.Debug("Inbound trunk: ", trunk.Name)
		found := false
		for inboundIdx := range trunks {
			it := trunks[inboundIdx]
			if it.Name == trunk.Name {
				foundInbounds[inboundIdx] = true
				found = true
				// Determine if we had changes that need to be reconciled
				updated := conversions.MergeSIPInboundTrunkInfos(
					trunk,
					conversions.MakeSIPInboundTrunkInfo(it))
				if trunk.String() != updated.String() {
					c.logger.Debugf("Inbound trunk configuration (%s) has changed. Updating.", trunk.Name)
					if c.shouldAssert() {
						return false
					}
					_, err := c.client.UpdateInboundSIPTrunk(updated)
					if err != nil {
						c.logger.Warnf("Failed to update inbound trunk (%s): %s", trunk.Name, err)
					}
				}
			}
		}
		if found == false {
			// Unidentified inbound trunk found
			if c.shouldAssert() {
				c.logger.Debug("Unexpected inbound trunk found: ", trunk.Name)
				return false
			}
			if c.shouldOverwrite() {
				c.logger.Debug("Removing Unexpected inbound trunk: ", trunk.Name)
				_, err := c.client.DeleteSIPTrunk(trunk.SipTrunkId)
				if err != nil {
					c.logger.Warnf("Failed to remove inbound trunk (%s): %s ", trunk.Name, err)
				}
			}
		}
	}
	for idx := range trunks {
		trunk := trunks[idx]
		if foundInbounds[idx] == false {
			if c.shouldAssert() {
				c.logger.Debug("Expected inbound trunk not found: ", trunk.Name)
				return false
			}
			_, err := c.client.CreateInboundSIPTrunk(
				conversions.MakeSIPInboundTrunkInfo(trunk))
			if err != nil {
				c.logger.Warnf("Failed to create inbound trunk (%s): %s ", trunk.Name, err)
			}
		}
	}
	return true
}

func (c *Controller) reconcileSIPOutboundTrunks(trunks []*types.SIPOutboundBridgeRequest) bool {
	if len(trunks) == 0 {
		return true
	}
	c.logger.Debug("Reconciling SIP outbound trunks ...")
	resp, err := c.client.ListSIPOutboundTrunk()
	if err != nil {
		c.logger.Warnf("failed to list SIP outbound trunks, skipping this round... %s", err)
		return false
	}
	c.logger.WithField("items", len(resp.Items)).Info("got SIP outbound trunks from server")

	foundOutbounds := make([]bool, len(trunks), len(trunks))
	for trunkIdx := range resp.Items {
		trunk := resp.Items[trunkIdx]
		c.logger.Debug("Outbound trunk: ", trunk.Name)
		found := false
		for inboundIdx := range trunks {
			it := trunks[inboundIdx]
			if it.Name == trunk.Name {
				foundOutbounds[inboundIdx] = true
				found = true
				// Determine if we had changes that need to be reconciled
				updated := conversions.MergeSIPOutboundTrunkInfos(
					trunk,
					conversions.MakeSIPOutboundTrunkInfo(it))
				if trunk.String() != updated.String() {
					c.logger.Debugf("Outbound trunk configuration (%s) has changed. Updating.", trunk.Name)
					if c.shouldAssert() {
						return false
					}
					_, err := c.client.UpdateOutboundSIPTrunk(updated)
					if err != nil {
						c.logger.Warnf("Failed to update outbound trunk (%s): %s", trunk.Name, err)
					}
				}
			}
		}
		if found == false {
			if c.shouldAssert() {
				c.logger.Debug("Unexpected outbound trunk found: ", trunk.Name)
				return false
			}
			if c.shouldOverwrite() {
				c.logger.Debug("Removing Unexpected outbound trunk: ", trunk.Name)
				_, err := c.client.DeleteSIPTrunk(trunk.SipTrunkId)
				if err != nil {
					c.logger.Warnf("Failed to remove outbound trunk (%s): %s ", trunk.Name, err)
				}
			}
		}
	}
	for idx := range trunks {
		trunk := trunks[idx]
		if foundOutbounds[idx] == false {
			if c.shouldAssert() {
				c.logger.Debug("Expected outbound trunk not found: ", trunk.Name)
				return false
			}
			_, err := c.client.CreateOutboundSIPTrunk(
				conversions.MakeSIPOutboundTrunkInfo(trunk))
			if err != nil {
				c.logger.Warnf("Failed to create inbound trunk (%s): %s ", trunk.Name, err)
			}
		}
	}
	return true
}

func (c *Controller) reconcileSIPDispatchRules(rules []*types.SIPDispatchRule) bool {
	if len(rules) == 0 {
		return true
	}
	c.logger.Debug("Reconciling SIP Dispatch Rules ...")
	resp, err := c.client.ListSIPDispatchRules()
	if err != nil {
		c.logger.Warnf("failed to list SIP dispatch rules, skipping this round... %s", err)
		return false
	}
	c.logger.WithField("items", len(resp.Items)).Info("got SIP Dispatch Rules from server")

	foundRules := make([]bool, len(rules), len(rules))
	for idx := range resp.Items {
		rule := resp.Items[idx]
		c.logger.Debug("Dispatch Rule: ", rule.Name)
		found := false
		for ruleIdx := range rules {
			r := rules[ruleIdx]
			if r.Name == rule.Name {
				foundRules[ruleIdx] = true
				found = true
				// Determine if we had changes that need to be reconciled
				updated := conversions.MergeSIPDispatchRules(
					rule,
					conversions.MakeSIPDispatchRule(r))
				if rule.String() != updated.String() {
					c.logger.Debugf("Dispatch Rule (%s) has changed. Updating.", rule.Name)
					if c.shouldAssert() {
						return false
					}
					_, err := c.client.UpdateSIPDispatchRule(updated)
					if err != nil {
						c.logger.Warnf("Failed to update Dispatch Rule (%s): %s", rule.Name, err)
					}
				}
			}
		}
		if found == false {
			if c.shouldAssert() {
				c.logger.Debug("Unexpected SIP Dispatch Rule found: ", rule.Name)
				return false
			}
			if c.shouldOverwrite() {
				c.logger.Debug("Removing Unexpected Dispatch Rule: ", rule.Name)
				_, err := c.client.DeleteSIPDispatchRule(rule.SipDispatchRuleId)
				if err != nil {
					c.logger.Warnf("Failed to remove Dispatch Rule (%s): %s ", rule.Name, err)
				}
			}
		}
	}
	for idx := range rules {
		rule := rules[idx]
		if foundRules[idx] == false {
			if c.shouldAssert() {
				c.logger.Debug("Expected Dispatch Rule not found: ", rule.Name)
				return false
			}
			_, err := c.client.CreateSIPDispatchRule(
				conversions.MakeSIPDispatchRule(rule))
			if err != nil {
				c.logger.Warnf("Failed to create Dispatch Rule (%s): %s ", rule.Name, err)
			}
		}
	}
	return true
}
