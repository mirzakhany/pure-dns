package puredns

import (
	"fmt"
	"log"

	"github.com/mirzakhany/pure_dns/internal/pkg/config"
	"github.com/mirzakhany/pure_dns/internal/pkg/dns"
	"golang.org/x/sync/errgroup"
)

// version params
// this variables will set with -ldflags on build time
var (
	LongHash   string
	CommitDate string
	BuildDate  string
)

// InitPureDNS InitPureDns application
func InitPureDNS(appName string, configPath string) error {
	// print version data to stdout
	printVersion(appName)

	// load config data
	config, err := config.LoadConf(appName, configPath)
	if err != nil {
		log.Printf("Load config file error: '%v'", err)
		return err
	}

	var g errgroup.Group

	g.Go(func() error {
		return dns.StartDNSServer()
	})

	fmt.Println(fmt.Sprintf("DNS Server Start On %s:%d", config.Server.Address, config.Server.Port))
	if err = g.Wait(); err != nil {
		log.Fatal(err)
	}

	return err
}

func printVersion(appName string) {
	fmt.Printf("===\nApp:%s\nCommitDate:%s\nBuildDate:%s\nLongHash:%s\n===\n",
		appName, CommitDate, BuildDate, LongHash)
}
