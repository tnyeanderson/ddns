package cmd

import (
	"log/slog"
	"os"
	"regexp"

	ddns "github.com/tnyeanderson/ddns/pkg"
	"gopkg.in/yaml.v3"
)

const (
	// EnvServerConfigFile is a path to a YAML config file for the server.
	EnvConfigFile = "DDNS_SERVER_CONFIG_FILE"

	EnvAPIServer          = "DDNS_API_SERVER"
	EnvAPIKey             = "DDNS_API_KEY"
	EnvServerAPIKey       = "DDNS_SERVER_API_KEY"
	EnvServerAPIKeyRegex  = "DDNS_SERVER_API_KEY_REGEX"
	EnvServerHostsFile    = "DDNS_SERVER_HOSTS_FILE"
	EnvServerHTTPListener = "DDNS_SERVER_HTTP_LISTENER"
	EnvServerDNSListener  = "DDNS_SERVER_DNS_LISTENER"
)

type Config struct {
	Agent  *ddns.Agent
	Server *ddns.Server
}

func (c *Config) Init() error {
	if c.Agent == nil {
		c.Agent = &ddns.Agent{}
	}

	if c.Server == nil {
		c.Server = &ddns.Server{}
	}

	// Read YAML config if it exists
	if v := os.Getenv(EnvConfigFile); v != "" {
		b, err := os.ReadFile(v)
		if err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
		if err := yaml.Unmarshal(b, c); err != nil {
			slog.Error(err.Error())
			os.Exit(1)
		}
	}

	// Overwrite config values with env vars, if set
	c.fromEnv()

	return nil
}

func (c *Config) fromEnv() {
	if v := os.Getenv(EnvServerHTTPListener); v != "" {
		c.Server.HTTPListener = v
	}

	if v := os.Getenv(EnvServerDNSListener); v != "" {
		c.Server.DNSListener = v
	}

	if key := os.Getenv(EnvServerAPIKey); key != "" {
		var r *regexp.Regexp
		if pattern := os.Getenv(EnvServerAPIKeyRegex); pattern != "" {
			r = regexp.MustCompile(pattern)
		}
		c.Server.Allow(key, r)
	}

	if v := os.Getenv(EnvServerHostsFile); v != "" {
		c.Server.HostsFile = v
	}

	if v := os.Getenv(EnvAPIServer); v != "" {
		c.Agent.ServerAddress = v
	}

	if v := os.Getenv(EnvAPIKey); v != "" {
		c.Agent.APIKey = v
	}
}
