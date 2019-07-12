package config

import (
	"strings"

	"bytes"
	"fmt"
	"io/ioutil"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

const (
	configName = "config"
)

var config ConfYaml

var defaultConf = []byte(`
server:
  port: 53
  address: 0.0.0.0

domains:
  - name: pure-dns.local
    host: localhost
    listenPort: 80
	targetPort: 8090
	
log:
	format: "string" # string or json
	access_log: "stdout" # stdout: output to console, or define log path like "log/access_log"
	access_level: "debug"
	error_log: "stderr" # stderr: output to console, or define log path like "log/error_log"
	error_level: "error"  
`)

// ConfYaml is config struct
type ConfYaml struct {
	Server  SectionServer    `yml:"server"`
	Domains []SectionDomains `yml:"domains"`
	Log     SectionLog       `yml:"log"`
}

// SectionServer is sub section of config
type SectionServer struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
}

// SectionDomains is sub section of config
type SectionDomains struct {
	Name       string `yaml:"name"`
	Host       string `yaml:"host"`
	listenPort int    `yaml:"listenPort"`
	targetPort int    `yaml:"targetPort"`
}

// SectionLog is sub section of config.
type SectionLog struct {
	Format      string `yaml:"format"`
	AccessLog   string `yaml:"access_log"`
	AccessLevel string `yaml:"access_level"`
	ErrorLog    string `yaml:"error_log"`
	ErrorLevel  string `yaml:"error_level"`
}

// Get return confYaml
func Get() ConfYaml {
	return config
}

// LoadConf load the config settings
func LoadConf(prefix string, configPath string) (ConfYaml, error) {

	var conf ConfYaml
	lowerPrefix := strings.ToLower(prefix)

	viper.SetConfigType("yaml")
	viper.AutomaticEnv()       // read in environment variables that match
	viper.SetEnvPrefix(prefix) // will be uppercased automatically
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if configPath != "" {
		content, err := ioutil.ReadFile(configPath)
		if err != nil {
			return conf, err
		}
		if err := viper.ReadConfig(bytes.NewBuffer(content)); err != nil {
			return conf, err
		}
	} else {

		viper.AddConfigPath(".")
		viper.AddConfigPath("/etc/" + lowerPrefix + "/")
		viper.AddConfigPath("$HOME/." + lowerPrefix)
		viper.SetConfigName(configName)

		// If a config file is found, read it in.
		if err := viper.ReadInConfig(); err == nil {
			fmt.Println("Using config file:", viper.ConfigFileUsed())

			viper.WatchConfig()
			viper.OnConfigChange(func(e fsnotify.Event) {
				fmt.Println("Config file changed:", e.Name)
			})

		} else {
			fmt.Println("load default config ...")
			// load default config
			if err := viper.ReadConfig(bytes.NewBuffer(defaultConf)); err != nil {
				return conf, err
			}
		}
	}

	var err error

	conf.Server.Address = viper.GetString("server.address")
	conf.Server.Port = viper.GetInt("server.port")

	err = viper.UnmarshalKey("domains", &conf.Domains)
	if err != nil {
		return conf, err
	}
	conf.Log.Format = viper.GetString("log.format")
	conf.Log.AccessLog = viper.GetString("log.access_log")
	conf.Log.AccessLevel = viper.GetString("log.access_level")
	conf.Log.ErrorLog = viper.GetString("log.error_log")
	conf.Log.ErrorLevel = viper.GetString("log.error_level")

	config = conf

	return conf, err
}
