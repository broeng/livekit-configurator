package livekit

import (
	"context"
	"time"

	"github.com/broeng/livekit-configurator/internal/config"
	auth "github.com/livekit/protocol/auth"
	"github.com/sirupsen/logrus"
	//lks "github.com/livekit/livekit-server"
	proto "github.com/livekit/protocol/livekit"
	livekit "github.com/livekit/server-sdk-go/v2"
)

type LiveKitClient struct {
	logger    logrus.FieldLogger
	context   context.Context
	serverUrl string
	apiKey    string
	apiSecret string
	sipClient *livekit.SIPClient
}

type grantsKey struct{}

type grantsValue struct {
	claims *auth.ClaimGrants
	apiKey string
}

func New(logger logrus.FieldLogger, context context.Context, config *config.Config) *LiveKitClient {
	sipClient := livekit.NewSIPClient(config.ServerUrl, config.ApiKey, config.ApiSecret)
	lkc := &LiveKitClient{
		logger:    logger.WithField("component", "livekit-client"),
		context:   context,
		serverUrl: config.ServerUrl,
		apiKey:    config.ApiKey,
		apiSecret: config.ApiSecret,
		sipClient: sipClient,
	}
	return lkc
}

func (c *LiveKitClient) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(c.context, 30*time.Second)
}

func (c *LiveKitClient) newContextWithGrants(grants *auth.ClaimGrants) (context.Context, context.CancelFunc) {
	base, cancel := c.newContext()
	ctx := context.WithValue(
		base,
		grantsKey{},
		&grantsValue{
			claims: grants,
			apiKey: c.apiKey,
		},
	)
	return ctx, cancel
}

func (c *LiveKitClient) newContextWithSIPGrants() (context.Context, context.CancelFunc) {
	return c.newContextWithGrants(&auth.ClaimGrants{
		SIP: &auth.SIPGrant{
			Admin: true,
			Call:  false,
		},
	})
}

func (c *LiveKitClient) TestConnection() bool {
	_, err := c.ListSIPInboundTrunk()
	return err == nil
}

func (c *LiveKitClient) ListSIPInboundTrunk() (*proto.ListSIPInboundTrunkResponse, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.ListSIPInboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.ListSIPInboundTrunkRequest{})
}

func (c *LiveKitClient) ListSIPOutboundTrunk() (*proto.ListSIPOutboundTrunkResponse, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.ListSIPOutboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.ListSIPOutboundTrunkRequest{})
}

func (c *LiveKitClient) ListSIPDispatchRules() (*proto.ListSIPDispatchRuleResponse, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.ListSIPDispatchRule(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.ListSIPDispatchRuleRequest{})
}

func (c *LiveKitClient) DeleteSIPTrunk(trunkId string) (*proto.SIPTrunkInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.DeleteSIPTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.DeleteSIPTrunkRequest{
			SipTrunkId: trunkId,
		})
}

func (c *LiveKitClient) DeleteSIPDispatchRule(ruleId string) (*proto.SIPDispatchRuleInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.DeleteSIPDispatchRule(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.DeleteSIPDispatchRuleRequest{
			SipDispatchRuleId: ruleId,
		})
}

func (c *LiveKitClient) CreateInboundSIPTrunk(info *proto.SIPInboundTrunkInfo) (*proto.SIPInboundTrunkInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.CreateSIPInboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.CreateSIPInboundTrunkRequest{
			Trunk: info,
		})
}

func (c *LiveKitClient) CreateOutboundSIPTrunk(info *proto.SIPOutboundTrunkInfo) (*proto.SIPOutboundTrunkInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.CreateSIPOutboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.CreateSIPOutboundTrunkRequest{
			Trunk: info,
		})
}

func (c *LiveKitClient) CreateSIPDispatchRule(rule *proto.SIPDispatchRuleInfo) (*proto.SIPDispatchRuleInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.CreateSIPDispatchRule(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.CreateSIPDispatchRuleRequest{
			DispatchRule: rule,
		})
}

func (c *LiveKitClient) UpdateInboundSIPTrunk(info *proto.SIPInboundTrunkInfo) (*proto.SIPInboundTrunkInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.UpdateSIPInboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.UpdateSIPInboundTrunkRequest{
			SipTrunkId: info.SipTrunkId,
			Action: &proto.UpdateSIPInboundTrunkRequest_Replace {
				Replace: info,
			},
		})
}

func (c *LiveKitClient) UpdateOutboundSIPTrunk(info *proto.SIPOutboundTrunkInfo) (*proto.SIPOutboundTrunkInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.UpdateSIPOutboundTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.UpdateSIPOutboundTrunkRequest{
			SipTrunkId: info.SipTrunkId,
			Action: &proto.UpdateSIPOutboundTrunkRequest_Replace {
				Replace: info,
			},
		})
}

func (c *LiveKitClient) UpdateSIPDispatchRule(rule *proto.SIPDispatchRuleInfo) (*proto.SIPDispatchRuleInfo, error) {
	ctx, cancel := c.newContextWithSIPGrants()
	defer cancel()
	return c.sipClient.UpdateSIPDispatchRule(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue{
				claims: &auth.ClaimGrants{},
				apiKey: c.apiKey,
			},
		),
		&proto.UpdateSIPDispatchRuleRequest{
			SipDispatchRuleId: rule.SipDispatchRuleId,
			Action: &proto.UpdateSIPDispatchRuleRequest_Replace {
				Replace: rule,
			},
		})
}
