package main

import (
	"flag"

	puredns "github.com/mirzakhany/pure_dns/internal/app/pure_dns"
)

const appName = "pure_dns"

func main() {
	// tcp address to start on. default localhost:53
	configPath := flag.String("config_path", "", "config Yaml file")
	flag.Parse()
	// init hub application with name and address flags
	puredns.InitPureDNS(appName, *configPath)
}
