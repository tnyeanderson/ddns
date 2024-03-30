package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"regexp"
	"strings"

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

func getEnvDocs(prefix string) string {
	// Not a map to preserve deterministic order
	docs := [][]string{
		[]string{EnvConfigFile, "Path to a YAML configuration file for DDNS. See the cmd.Config struct for more info."},
		[]string{EnvAPIServer, "(Agent) The scheme/host/port of the DDNS API server, not including the /api base path."},
		[]string{EnvAPIKey, "(Agent) The API key used to authenticate to the DDNS API."},
		[]string{EnvServerAPIKey, "Adds this API key to the AllowedAPIKeys map for the server."},
		[]string{EnvServerAPIKeyRegex, fmt.Sprintf("Used as the value for the %s key in the Server.AllowedAPIKeys map.", EnvServerAPIKey)},
		[]string{EnvServerHostsFile, "Path to the hosts file used by the server. See Server.HostsFile for more info."},
		[]string{EnvServerHTTPListener, "The TCP listener address for the HTTP server."},
		[]string{EnvServerDNSListener, "The TCP listener address for the DNS server."},
	}
	s := &strings.Builder{}
	for _, v := range docs {
		fmt.Fprintf(s, "%s%s\n", prefix, v[0])
		fmt.Fprintf(s, "%s  %s\n\n", prefix, v[1])
	}
	return s.String()
}
