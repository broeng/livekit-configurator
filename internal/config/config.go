package config

type ConfigT struct {
	LogLevel    string `id:"log-level" short:"l" desc:"verbosity level for logs" default:"info"`
	RunOnce     bool `id:"run-once" desc:"run once and exit" default:"false"`

	ServerUrl   string `id:"server-url" short:"s" desc:"server url for LiveKit server"`
	ApiKey      string `id:"api-key" desc:"name of API key to use for LiveKit server"`
	ApiSecret   string `id:"api-secret" desc:"secret key for specified API key"`

	ConfigPath  string `id:"config" desc:"path to livekit configuration definition"`

	ReconciliationFrequency int `id:"delay" desc:"delay between reconciliation attempts, in seconds" default:"10"`

	Version     bool `id:"version" desc:"show version and quit" opts:"hidden"`
}

var Config ConfigT
