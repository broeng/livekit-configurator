package livekit

import (
	"context"
	"time"

	"github.com/broeng/livekit-configurator/internal/config"
	"github.com/sirupsen/logrus"
	auth "github.com/livekit/protocol/auth"
	//lks "github.com/livekit/livekit-server"
	livekit "github.com/livekit/server-sdk-go/v2"
	proto "github.com/livekit/protocol/livekit"
)

type LiveKitClient struct {
	logger     logrus.FieldLogger
	context    context.Context
	serverUrl  string
	apiKey     string
	apiSecret  string
	sipClient  *livekit.SIPClient
}

type grantsKey struct {}

type grantsValue struct {
	claims *auth.ClaimGrants
	apiKey string
}

func New(logger logrus.FieldLogger, context context.Context, config config.ConfigT) *LiveKitClient {
	logger.Info("Using api key: ", config.ApiKey)
	logger.Info("Using api secret: ", config.ApiSecret)
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
		&grantsValue {
			claims: grants,
			apiKey: c.apiKey,
		},
	)
	return ctx, cancel
}

func (c *LiveKitClient) ListSIPTrunks() (*proto.ListSIPTrunkResponse, error) {
	//at := auth.NewAccessToken(c.apiKey, c.apiSecret)
	//grant := &auth.SIPGrant{
	//	Admin: true,
	//	Call: false,
	//}
	//at.AddSIPGrant(grant).SetIdentity(c.apiKey).SetValidFor(time.Hour)
	//jwt, err := at.ToJWT()
	//if err != nil {
	//	return nil, err
	//}
	ctx, cancel := c.newContextWithGrants(&auth.ClaimGrants{
		SIP: &auth.SIPGrant{
			Admin: true,
			Call: false,
		},
	})
	defer cancel()
	return c.sipClient.ListSIPTrunk(
		context.WithValue(
			ctx,
			grantsKey{},
			&grantsValue {
				claims: &auth.ClaimGrants{
				},
				apiKey: c.apiKey,
			},
		),
		&proto.ListSIPTrunkRequest {})
}
