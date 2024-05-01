package main

import (
	"github.com/tnyeanderson/ddns/cmd"
)

func main() {
	cmd.SetVersion(version())
	cmd.Execute()
}
