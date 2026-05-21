package conversions

import (
	"time"

	types "github.com/broeng/livekit-configurator/internal/types"
	proto "github.com/livekit/protocol/livekit"
	durationpb "google.golang.org/protobuf/types/known/durationpb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func MakeSIPInboundTrunkInfo(desc *types.SIPInboundBridgeRequest) *proto.SIPInboundTrunkInfo {
	emptyMap := map[string]string{}
	return &proto.SIPInboundTrunkInfo{
		Name:                desc.Name,
		Numbers:             desc.Numbers,
		AllowedAddresses:    desc.AllowedAddresses,
		AllowedNumbers:      desc.AllowedNumbers,
		AuthUsername:        desc.Username,
		AuthPassword:        desc.Password,
		Headers:             emptyMap,
		HeadersToAttributes: emptyMap,
		AttributesToHeaders: emptyMap,
		IncludeHeaders:      proto.SIPHeaderOptions_SIP_NO_HEADERS,
		RingingTimeout:      durationpb.New(120 * time.Second),
		MaxCallDuration:     durationpb.New(24 * time.Hour),
		//KrispEnabled: false, // TODO
		MediaEncryption: proto.SIPMediaEncryption_SIP_MEDIA_ENCRYPT_ALLOW,
		CreatedAt:       tspb.Now(),
		UpdatedAt:       tspb.Now(),
	}
}

func MergeSIPInboundTrunkInfos(base *proto.SIPInboundTrunkInfo, prio *proto.SIPInboundTrunkInfo) *proto.SIPInboundTrunkInfo {
	return &proto.SIPInboundTrunkInfo{
		SipTrunkId:          base.SipTrunkId,
		Name:                prio.Name,
		Numbers:             prio.Numbers,
		AllowedAddresses:    prio.AllowedAddresses,
		AllowedNumbers:      prio.AllowedNumbers,
		AuthUsername:        prio.AuthUsername,
		AuthPassword:        prio.AuthPassword,
		Headers:             prio.Headers,
		HeadersToAttributes: prio.HeadersToAttributes,
		AttributesToHeaders: prio.AttributesToHeaders,
		IncludeHeaders:      prio.IncludeHeaders,
		RingingTimeout:      prio.RingingTimeout,
		MaxCallDuration:     prio.MaxCallDuration,
		KrispEnabled:        prio.KrispEnabled,
		MediaEncryption:     prio.MediaEncryption,
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
	}
}

func MakeSIPOutboundTrunkInfo(desc *types.SIPOutboundBridgeRequest) *proto.SIPOutboundTrunkInfo {
	emptyMap := map[string]string{}
	return &proto.SIPOutboundTrunkInfo{
		Name:                desc.Name,
		Address:             desc.Address,
		DestinationCountry:  desc.DestinationCountry,
		Transport:           proto.SIPTransport_SIP_TRANSPORT_AUTO,
		Numbers:             desc.Numbers,
		AuthUsername:        desc.Username,
		AuthPassword:        desc.Password,
		Headers:             emptyMap,
		HeadersToAttributes: emptyMap,
		AttributesToHeaders: emptyMap,
		IncludeHeaders:      proto.SIPHeaderOptions_SIP_NO_HEADERS,
		MediaEncryption:     proto.SIPMediaEncryption_SIP_MEDIA_ENCRYPT_ALLOW,
		FromHost:            desc.FromHost,
		CreatedAt:           tspb.Now(),
		UpdatedAt:           tspb.Now(),
	}
}

func MergeSIPOutboundTrunkInfos(base *proto.SIPOutboundTrunkInfo, prio *proto.SIPOutboundTrunkInfo) *proto.SIPOutboundTrunkInfo {
	return &proto.SIPOutboundTrunkInfo{
		SipTrunkId:          base.SipTrunkId,
		Name:                prio.Name,
		Address:             prio.Address,
		DestinationCountry:  prio.DestinationCountry,
		Transport:           prio.Transport,
		Numbers:             prio.Numbers,
		AuthUsername:        prio.AuthUsername,
		AuthPassword:        prio.AuthPassword,
		Headers:             prio.Headers,
		HeadersToAttributes: prio.HeadersToAttributes,
		AttributesToHeaders: prio.AttributesToHeaders,
		IncludeHeaders:      prio.IncludeHeaders,
		MediaEncryption:     prio.MediaEncryption,
		FromHost:            prio.FromHost,
		CreatedAt:           base.CreatedAt,
		UpdatedAt:           base.UpdatedAt,
	}
}

func makeDispatchRule(desc *types.DispatchRule) *proto.SIPDispatchRule {
	if desc.DispatchRuleDirect != nil {
		return &proto.SIPDispatchRule{
			Rule: &proto.SIPDispatchRule_DispatchRuleDirect{
				DispatchRuleDirect: &proto.SIPDispatchRuleDirect{
					RoomName: desc.DispatchRuleDirect.RoomName,
					Pin:      desc.DispatchRuleDirect.Pin,
				},
			},
		}
	} else if desc.DispatchRuleCallee != nil {
		return &proto.SIPDispatchRule{
			Rule: &proto.SIPDispatchRule_DispatchRuleCallee{
				DispatchRuleCallee: &proto.SIPDispatchRuleCallee{
					RoomPrefix: desc.DispatchRuleCallee.RoomPrefix,
					Pin:        desc.DispatchRuleCallee.Pin,
					Randomize:  desc.DispatchRuleCallee.Randomize,
				},
			},
		}
	} else if desc.DispatchRuleIndividual != nil {
		return &proto.SIPDispatchRule{
			Rule: &proto.SIPDispatchRule_DispatchRuleIndividual{
				DispatchRuleIndividual: &proto.SIPDispatchRuleIndividual{
					RoomPrefix:   desc.DispatchRuleIndividual.RoomPrefix,
					Pin:          desc.DispatchRuleIndividual.Pin,
					NoRandomness: desc.DispatchRuleIndividual.NoRandomness,
				},
			},
		}
	} else {
		return &proto.SIPDispatchRule{
			Rule: &proto.SIPDispatchRule_DispatchRuleDirect{
				DispatchRuleDirect: &proto.SIPDispatchRuleDirect{
					RoomName: "default",
				},
			},
		}
	}
}

func MakeSIPDispatchRule(desc *types.SIPDispatchRule) *proto.SIPDispatchRuleInfo {
	emptyMap := map[string]string{}
	return &proto.SIPDispatchRuleInfo{
		Name:            desc.Name,
		Rule:            makeDispatchRule(&desc.Rule),
		Numbers:         desc.Numbers,
		HidePhoneNumber: desc.HidePhoneNumber,
		InboundNumbers:  desc.InboundNumbers,
		Attributes:      emptyMap,
		//RoomConfig:        // not needed?
		//KrispEnabled:        false, // TODO
		MediaEncryption: proto.SIPMediaEncryption_SIP_MEDIA_ENCRYPT_ALLOW,
		CreatedAt:       tspb.Now(),
		UpdatedAt:       tspb.Now(),
	}
}

func MergeSIPDispatchRules(base *proto.SIPDispatchRuleInfo, prio *proto.SIPDispatchRuleInfo) *proto.SIPDispatchRuleInfo {
	return &proto.SIPDispatchRuleInfo{
		SipDispatchRuleId: base.SipDispatchRuleId,
		Name:              prio.Name,
		Rule:              prio.Rule,
		Numbers:           prio.Numbers,
		HidePhoneNumber:   prio.HidePhoneNumber,
		InboundNumbers:    prio.InboundNumbers,
		Attributes:        prio.Attributes,
		RoomConfig:        prio.RoomConfig,
		KrispEnabled:      prio.KrispEnabled,
		MediaEncryption:   prio.MediaEncryption,
		CreatedAt:         base.CreatedAt,
		UpdatedAt:         base.UpdatedAt,
	}
}
