package types

type SIPConfig struct {
	BridgeRequest   *SIPBridgeRequest `json:"bridgeRequest",omitempty`
	DispatchRule    *SIPDispatchRule  `json:"dispatchRule",omitempty`
}

type LiveKitConfigDefinition struct {
	SIPConfig   *SIPConfig  `json:"sip",omitempty`
}
