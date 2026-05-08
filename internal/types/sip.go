package types

type SIPBridgeRequest struct {
	// Name used for identifying this SIP bridge
	Name                  string `json:"name"`
	// Array of IP Address or CIDRs where SIP INVITEs will be accepted from
	InboundAddresses      []string `json:"inbound_addresses"`
	// IP Address that SIP INVITEs will be sent too
	OutboundAddress       string `json:"outbound_address"`
	// When making an outbound call on this SIP Trunk what Phone Number should be used
	OutboundNumber        string `json:"outbound_number"`
	// Phone numbers this SIP Trunk will serve. If Empty it will serve all incoming calls,
	InboundNumbersRegex   string `json:"inbound_numbers_regex"`
	// Username for Authentication of inbound calls, no Authentication if empty,
	InboundUsername       string `json:"inbound_username",omitempty`
	// Password for Authentication of inbound calls, no Authentication if empty,
	InboundPassword       string `json:"inbound_password",omitempty`
	// Password file for Authentication of inbound calls, if file exists, the
	// content of the file will be used as the password. Takes precedence over
	// InboundPassword.
	InboundPasswordFile   string `json:"inbound_password_file",omitempty`
	// Username for Authentication of outbound calls, no Authentication if empty,
	OutboundUsername       string `json:"outbound_username",omitempty`
	// Password for Authentication of outbound calls, no Authentication if empty,
	OutboundPassword       string `json:"outbound_password",omitempty`
	// Password file for Authentication of outbound calls, if file exists, the
	// content of the file will be used as the password. Takes precedence over
	// OutboundPassword.
	OutboundPasswordFile   string `json:"outbound_password_file",omitempty`
}

type DispatchDirectRule struct {
	RoomName string `json:"roomName"`
}

type DispatchRule struct {
	// Only one rule supported for now
	DispatchRuleDirect DispatchDirectRule `json:"dispatchRuleDirect"`
}

type SIPDispatchRule struct {
	// What rule to use to dispatch this call
	Rule             DispatchRule `json:"rule"`
	// Array of SIP Trunk IDs that are accepted for this rule. If empty all Trunks
	TrunkIds         []string `json:"trunk_ids"`
	// If true hide the phone number when joining the LiveKit room
	HidePhoneNumber  bool `json:"hide_phone_number"`
}
