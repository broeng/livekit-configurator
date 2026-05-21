package config

import (
	"fmt"
	"github.com/stevenroose/gonfig"
)

type Config struct {
	LogLevel string `id:"log-level" short:"l" desc:"verbosity level for logs" default:"info"`
	RunOnce  bool   `id:"run-once" desc:"run once and exit" default:"false"`

	ServerUrl string `id:"server-url" short:"s" desc:"server url for LiveKit server"`
	ApiKey    string `id:"api-key" desc:"name of API key to use for LiveKit server"`
	ApiSecret string `id:"api-secret" desc:"secret key for specified API key"`

	ConfigPath string `id:"config" short:"c" desc:"path to livekit configuration definition"`

	ReconciliationMode      *Mode `id:"mode" desc:"reconciliation mode (one of: assert, merge, overwrite)"`
	ReconciliationFrequency int   `id:"delay" desc:"delay between reconciliation attempts, in seconds" default:"10"`

	HealthListenPort int `id:"port" desc:"listen port for http health service" default:"8080"`

	Version bool `id:"version" desc:"show version and quit" opts:"hidden"`
}

func LoadConfig(envPrefix string) (*Config, error) {
	var config Config
	if err := gonfig.Load(&config, gonfig.Conf{
		EnvPrefix: envPrefix,
		FlagIgnoreUnknown: false,
	}); err != nil {
		return nil, fmt.Errorf("could not parse options: %s", err)
	}
	if !config.Version {
		if len(config.ConfigPath) == 0 {
			return nil, fmt.Errorf("--config/%sCONFIG parameter missing", envPrefix)
		}
		if len(config.ServerUrl) == 0 {
			return nil, fmt.Errorf("--server-url/%sSERVER_URL parameter missing", envPrefix)
		}
		if len(config.ApiKey) == 0 {
			return nil, fmt.Errorf("--api-key/%sAPI_KEY parameter missing", envPrefix)
		}
		if len(config.ApiSecret) == 0 {
			return nil, fmt.Errorf("--api-secret/%sAPI_SECRET parameter missing", envPrefix)
		}
	}
	return &config, nil
}
