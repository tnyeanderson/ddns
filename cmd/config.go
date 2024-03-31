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
	// EnvConfigFile is the path to a YAML config file for the server.
	EnvConfigFile = "DDNS_CONFIG_FILE"

	EnvAPIServer = "DDNS_API_SERVER" // sets [Agent.ServerAddress]
	EnvAPIKey    = "DDNS_API_KEY"    // sets [Agent.APIKey]

	EnvServerAPIKey       = "DDNS_SERVER_API_KEY"       // sets a key in [Server.AllowedAPIKeys]
	EnvServerAPIKeyRegex  = "DDNS_SERVER_API_KEY_REGEX" // sets the value for the [EnvServerAPIKey] key in [Server.AllowedAPIKeys]
	EnvServerHostsFile    = "DDNS_SERVER_HOSTS_FILE"    // sets [Server.HostsFile]
	EnvServerHTTPListener = "DDNS_SERVER_HTTP_LISTENER" // sets [Server.HTTPListener]
	EnvServerDNSListener  = "DDNS_SERVER_DNS_LISTENER"  // sets [Server.DNSListener]
)

// Config contains the configuration used by the ddns CLI. It is also the data
// structure used when unmarshaling the YAML config file.
type Config struct {
	Agent  *ddns.Agent
	Server *ddns.Server
}

// Init tries to set the values in [c], first using a YAML config file (if
// provided), then using environment variables.
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
		[]string{EnvConfigFile, `Path to a YAML configuration file for DDNS. See the cmd.Config struct for more info.`},
		[]string{EnvAPIServer, fmt.Sprintf(`(Agent) The scheme/host/port of the DDNS API server, not including the /api base path (default: "%s").`, ddns.DefaultServerAddress)},
		[]string{EnvAPIKey, `(Agent) The API key used to authenticate to the DDNS API.`},
		[]string{EnvServerAPIKey, `This API key will be allowed by the server.`},
		[]string{EnvServerAPIKeyRegex, fmt.Sprintf(`The regex domain matcher for %s.`, EnvServerAPIKey)},
		[]string{EnvServerHostsFile, `Path to the hosts file used by the server. See Server.HostsFile for more info.`},
		[]string{EnvServerHTTPListener, fmt.Sprintf(`The TCP listener address for the HTTP server (default: "%s").`, ddns.DefaultHTTPListener)},
		[]string{EnvServerDNSListener, fmt.Sprintf(`The TCP listener address for the DNS server (default: "%s").`, ddns.DefaultDNSListener)},
	}
	s := &strings.Builder{}
	for _, v := range docs {
		fmt.Fprintf(s, "%s%s\n", prefix, v[0])
		fmt.Fprintf(s, "%s  %s\n\n", prefix, v[1])
	}
	return s.String()
}
