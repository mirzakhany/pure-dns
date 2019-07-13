package config

import (
	"strings"

	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/spf13/viper"
)

const (
	configName = "config"
)

// Config hold application config
var Config ConfYaml

var defaultConf = []byte(`
server:
  port: 53
  address: 127.0.0.1

domains:
  - name: pure-dns.local
    host: localhost

log:
    format: "string" 
    access_log: "stdout" 
    access_level: "debug"
    error_log: "stderr"
    error_level: "error"
`)

// ConfYaml is config struct
type ConfYaml struct {
	Server    SectionServer    `yaml:"server"`
	Domains   []SectionDomains `yaml:"domains"`
	Log       SectionLog       `yaml:"log"`
	DomainMap map[string]string
}

// SectionServer is sub section of config
type SectionServer struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// SectionDomains is sub section of config
type SectionDomains struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Format      string `yaml:"format"`
	AccessLog   string `yaml:"access_log"`
	AccessLevel string `yaml:"access_level"`
	ErrorLog    string `yaml:"error_log"`
	ErrorLevel  string `yaml:"error_level"`
}

func convertToMap(domains []SectionDomains) map[string]string {
	var resule = make(map[string]string)
	for _, domain := range domains {
		resule[domain.Name] = domain.Host
	}
	return resule
}

func loadData() (ConfYaml, error) {

	var conf ConfYaml
	conf.Server.Address = viper.GetString("server.address")
	conf.Server.Port = viper.GetInt("server.port")

	err := viper.UnmarshalKey("domains", &conf.Domains)
	if err != nil {
		return conf, err
	}
	conf.Log.Format = viper.GetString("log.format")
	conf.Log.AccessLog = viper.GetString("log.access_log")
	conf.Log.AccessLevel = viper.GetString("log.access_level")
	conf.Log.ErrorLog = viper.GetString("log.error_log")
	conf.Log.ErrorLevel = viper.GetString("log.error_level")

	conf.DomainMap = convertToMap(conf.Domains)
	return conf, err
}

// LoadConf load the config settings
func LoadConf(prefix string, configPath string) (ConfYaml, error) {

	lowerPrefix := strings.ToLower(prefix)

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()       // read in environment variables that match
	viper.SetEnvPrefix(prefix) // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configPath != "" {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return Config, err
		}
		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return Config, err
		}
	} else {

		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/" + lowerPrefix + "/")
		viper.AddConfigPath("$HOME/." + lowerPrefix)
		viper.SetConfigName(configName)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		} else {
			fmt.Println("load default config ...")
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return Config, err
			}
		}
	}

	Config, err := loadData()
	return Config, err
}
