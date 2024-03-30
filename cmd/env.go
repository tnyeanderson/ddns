package cmd

import (
	_ "embed"
)

////go:embed env.go
//var thisFile []byte
//
const (
	// EnvConfigFile is a path to a YAML config file.
	EnvConfigFile = "DDNS_CONFIG_FILE"

	// EnvAPIServer sets [Agent.ServerAddress].
	EnvAPIServer = "DDNS_API_SERVER"

	// EnvAPIKey sets [Agent.APIKey].
	EnvAPIKey = "DDNS_API_KEY"

	// EnvServerAPIKey adds the key to [Server.AllowedAPIKeys].
	EnvServerAPIKey = "DDNS_SERVER_API_KEY"

	// EnvServerAPIKeyRegex is compiled with [regexp.MustCompile()] and used as
	// the value for the key [EnvServerAPIKey] in [Server.AllowedAPIKeys].
	EnvServerAPIKeyRegex = "DDNS_SERVER_API_KEY_REGEX"

	// EnvServerHostsFile is used as the value for [Server.HostsFile].  Values in
	// this file will override the contents of "domains" in the config file.
	EnvServerHostsFile = "DDNS_SERVER_HOSTS_FILE"

	// EnvServerHTTPListener is used as the value for [Server.HTTPListener].
	EnvServerHTTPListener = "DDNS_SERVER_HTTP_LISTENER"

	// EnvServerDNSListener is used as the value for [Server.DNSListener].
	EnvServerDNSListener = "DDNS_SERVER_DNS_LISTENER"
)

//
//func envVarDocs() error {
//	fset := token.NewFileSet()
//	file, err := parser.ParseFile(fset, "env.go", thisFile, parser.ParseComments+parser.SkipObjectResolution)
//	if err != nil {
//		return "", err
//	}
//
//	pkg, err := doc.NewFromFiles(fset, []*ast.File{file}, "github.com/tnyeanderson/ddns/cmd")
//	if err != nil {
//		return "", err
//	}
//
//	out := ""
//	for _, c := range pkg.Consts {
//		for _, n := range c.Names {
//			out += fmt.Sprintln(n)
//		}
//	}
//
//	//ast.Inspect(file, func(n ast.Node) bool {
//	//	if x, ok := n.(*ast.ValueSpec); ok {
//	//		for _, ident := range x.Names {
//	//			if strings.HasPrefix(ident.Name, "Env") {
//	//				fmt.Println(x.Doc.Text())
//	//				fmt.Println(ident.Name)
//	//			}
//	//		}
//	//	}
//	//	return true
//	//})
//
//	return out, nil
//}
