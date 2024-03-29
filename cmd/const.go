package cmd

const (
	// EnvAPIServer sets [Agent.ServerAddress].
	EnvAPIServer = "DDNS_API_SERVER"

	// EnvAPIKey sets [Agent.APIKey].
	EnvAPIKey = "DDNS_API_KEY"

	// EnvServerAPIKey adds the key to [Server.AllowedAPIKeys].
	EnvServerAPIKey = "DDNS_SERVER_API_KEY"

	// EnvServerAPIKeyRegex is compiled with [regexp.MustCompile()] and used as
	// the value for the key [EnvServerAPIKey] in [Server.AllowedAPIKeys].
	EnvServerAPIKeyRegex = "DDNS_SERVER_API_KEY_REGEX"

	// EnvServerConfigFile is a path to a YAML config file for the server.
	EnvServerConfigFile = "DDNS_SERVER_CONFIG_FILE"

	// EnvServerHostsFile is used as the value for [Server.HostsFile].  Values in
	// this file will override the contents of "domains" in the config file.
	EnvServerHostsFile = "DDNS_SERVER_HOSTS_FILE"

	// EnvServerHTTPListener is used as the value for [Server.HTTPListener].
	EnvServerHTTPListener = "DDNS_SERVER_HTTP_LISTENER"

	// EnvServerDNSListener is used as the value for [Server.DNSListener].
	EnvServerDNSListener = "DDNS_SERVER_DNS_LISTENER"
)
