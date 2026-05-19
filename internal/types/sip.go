package types

type SIPInboundBridgeRequest struct {
	// Name used for identifying this SIP bridge
	Name string `json:"name"`
	// Array of IP Address or CIDRs where SIP INVITEs will be accepted from
	AllowedAddresses []string `json:"allowed_addresses"`
	// Numbers associated with LiveKit SIP. The Trunk will only accept calls made to these numbers.
	// Creating multiple Trunks with different phone numbers allows having different rules for a single provider.
	Numbers []string `json:"numbers"`
	// Phone numbers this SIP Trunk will serve. If Empty it will serve all incoming calls,
	AllowedNumbers []string `json:"allowed_numbers",omitempty`
	// Username for Authentication of inbound calls, no Authentication if empty,
	Username string `json:"username",omitempty`
	// Password for Authentication of inbound calls, no Authentication if empty,
	Password string `json:"password",omitempty`
	// Password file for Authentication of inbound calls, if file exists, the
	// content of the file will be used as the password. Takes precedence over
	// InboundPassword.
	PasswordFile string `json:"password_file",omitempty`
}

type SIPOutboundBridgeRequest struct {
	// Name used for identifying this SIP bridge
	Name string `json:"name"`
	// IP Address that SIP INVITEs will be sent too
	Address string `json:"address"`
	// country where the call terminates as ISO 3166-1 alpha-2 (https://en.wikipedia.org/wiki/ISO_3166-1_alpha-2).
	// This will be used by the livekit infrastructure to route calls.
	DestinationCountry string `json:"destination_country",omitempty`
	// Numbers used to make the calls. Random one from this list will be selected.
	Numbers []string `json:"numbers"`
	// Optional custom hostname for the 'From' SIP header in outbound INVITEs.
	// When set, outbound calls from this trunk will use this host instead of the default project SIP domain.
	// Enables originating calls from custom domains.
	FromHost string `json:"from_host",omitempty`
	// Username for Authentication of outbound calls, no Authentication if empty,
	Username string `json:"username",omitempty`
	// Password for Authentication of outbound calls, no Authentication if empty,
	Password string `json:"password",omitempty`
	// Password file for Authentication of outbound calls, if file exists, the
	// content of the file will be used as the password. Takes precedence over
	// OutboundPassword.
	PasswordFile string `json:"password_file",omitempty`
}

type SIPTrunks struct {
	InboundTrunks  []*SIPInboundBridgeRequest  `json:"inbound",omitempty`
	OutboundTrunks []*SIPOutboundBridgeRequest `json:"outbound",omitempty`
}

type DispatchDirectRule struct {
	RoomName string `json:"roomName"`
	Pin string `json:"pin",omitempty`
}

type DispatchIndividualRule struct {
	RoomPrefix string `json:"roomPrefix",omitempty`
	Pin string `json:"pin",omitempty`
	NoRandomness bool `json:"no_randomness",omitempty`
}

type DispatchCalleeRule struct {
	RoomPrefix string `json:"roomPrefix",omitempty`
	Pin string `json:"pin",omitempty`
	Randomize bool `json:"randomize",omitempty`
}

type DispatchRule struct {
	DispatchRuleDirect *DispatchDirectRule `json:"dispatchRuleDirect",omitempty`
	DispatchRuleIndividual *DispatchIndividualRule `json:"dispatchRuleIndividual",omitempty`
	DispatchRuleCallee *DispatchCalleeRule `json:"dispatchRuleCallee",omitempty`
}

type SIPDispatchRule struct {
	// Human-readable name for the Dispatch Rule.
	Name string `json:"name"`
	// What rule to use to dispatch this call
	Rule DispatchRule `json:"rule"`
	// Array of SIP Trunk IDs that are accepted for this rule. If empty all Trunks
	TrunkIds []string `json:"trunk_ids"`
	// Dispatch Rule will only accept a call made to these numbers (if set).
	Numbers  []string `json:"numbers",omitempty`
	// Dispatch Rule will only accept a call made from these numbers (if set).
	InboundNumbers  []string `json:"inbound_numbers",omitempty`
	// If true hide the phone number when joining the LiveKit room
	HidePhoneNumber bool `json:"hide_phone_number"`
}
