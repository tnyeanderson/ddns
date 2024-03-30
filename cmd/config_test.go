package cmd

import (
	"net"
	"os"
	"regexp"
	"testing"

	"github.com/go-test/deep"
	ddns "github.com/tnyeanderson/ddns/pkg"
)

func TestYAMLConfig(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Fatal(err.Error())
	}
	os.Setenv(EnvConfigFile, "testdata/ddns.yaml")
	c := Config{}
	if err := c.Init(); err != nil {
		t.Fatal(err.Error())
	}
	expected := Config{
		Agent: &ddns.Agent{
			ServerAddress: "https://myserver.com:1234",
			APIKey:        "mysuperdupersecret",
		},
		Server: &ddns.Server{
			HostsFile: "/path/to/hostsfile",
			AllowedAPIKeys: map[string]*regexp.Regexp{
				"mysupersecretkey": regexp.MustCompile("^onlythishost.com$"),
			},
			HTTPListener: ":8888",
			DNSListener:  ":5333",
			Domains: ddns.Domains{
				"domain1.haha": net.ParseIP("4.3.2.1"),
			},
		},
	}

	if diff := deep.Equal(c, expected); diff != nil {
		t.Fatalf("unmarshaled YAML config is incorrect: %v", diff)
	}
}

func TestYAMLConfigEnvOverride(t *testing.T) {
	if err := clearEnv(); err != nil {
		t.Fatal(err.Error())
	}
	os.Setenv(EnvConfigFile, "testdata/ddns.yaml")

	// These env vars should override the values from the config file
	envVals := map[string]string{
		EnvAPIServer:          "http://serverfromenv.com",
		EnvAPIKey:             "apikeyfromenv",
		EnvServerAPIKey:       "allowedkeyfromenv",
		EnvServerAPIKeyRegex:  ".*",
		EnvServerHostsFile:    "/path/to/hostsfile/from/env",
		EnvServerHTTPListener: ":1111",
		EnvServerDNSListener:  ":9999",
	}

	for k, v := range envVals {
		os.Setenv(k, v)
	}

	c := Config{}
	if err := c.Init(); err != nil {
		t.Fatal(err.Error())
	}
	expected := Config{
		Agent: &ddns.Agent{
			ServerAddress: envVals[EnvAPIServer],
			APIKey:        envVals[EnvAPIKey],
		},
		Server: &ddns.Server{
			HostsFile: envVals[EnvServerHostsFile],
			AllowedAPIKeys: map[string]*regexp.Regexp{
				"mysupersecretkey":       regexp.MustCompile("^onlythishost.com$"),
				envVals[EnvServerAPIKey]: regexp.MustCompile(envVals[EnvServerAPIKeyRegex]),
			},
			HTTPListener: envVals[EnvServerHTTPListener],
			DNSListener:  envVals[EnvServerDNSListener],
			Domains: ddns.Domains{
				"domain1.haha": net.ParseIP("4.3.2.1"),
			},
		},
	}

	if diff := deep.Equal(c, expected); diff != nil {
		t.Fatalf("unmarshaled YAML config is incorrect: %v", diff)
	}
}

func clearEnv() error {
	all := []string{
		EnvConfigFile,
		EnvAPIServer,
		EnvAPIKey,
		EnvServerAPIKey,
		EnvServerAPIKeyRegex,
		EnvServerHostsFile,
		EnvServerHTTPListener,
		EnvServerDNSListener,
	}
	for _, e := range all {
		if err := os.Unsetenv(e); err != nil {
			return err
		}
	}

	return nil
}
