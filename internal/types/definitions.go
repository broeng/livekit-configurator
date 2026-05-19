package types

type SIPConfig struct {
	Trunks        *SIPTrunks         `json:"trunks",omitempty`
	DispatchRules []*SIPDispatchRule `json:"dispatchRules",omitempty`
}

type LiveKitConfigDefinition struct {
	SIPConfig *SIPConfig `json:"sip",omitempty`
}
