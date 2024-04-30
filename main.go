package main

import (
	"fmt"

	"github.com/tnyeanderson/ddns/cmd"
)

// These get populated by goreleaser
var (
	version = "dev"
	commit  = "none"
)

func main() {
	cmd.Execute(fmt.Sprintf("%s commit:%s", version, commit))
}
