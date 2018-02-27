package server

import (
	"fmt"
	"io/ioutil"
	"ngrok/log"
	"path"

	"github.com/kardianos/osext"
	yaml "gopkg.in/yaml.v1"
)

type Configuration struct {
	HttpAddr         string `yaml:"http_addr,omitempty"`
	HttpsAddr        string `yaml:"https_addr,omitempty"`
	TunnelAddr       string `yaml:"tunnel_addr,omitempty"`
	Domain           string `yaml:"domain,omitempty"`
	TlsCrt           string `yaml:"tls_crt,omitempty"`
	TlsKey           string `yaml:"tls_key,omitempty"`
	LogTo            string `yaml:"log_to,omitempty"`
	LogLevel         string `yaml:"log_level,omitempty"`
	HttpPulbishPort  string `yaml:"http_pulbish_port,omitempty"`
	HttpsPulbishPort string `yaml:"https_pulbish_port,omitempty"`
}

const (
	defaultHTTPAddr   = ":80"
	defaultHTTPSAddr  = ":443"
	defaultDomain     = "t.mgrok.cn"
	defaultLogto      = "stdout"
	defaultLogLevel   = "DEBUG"
	defaultTunnelAddr = ":4443"
)

func LoadConfiguration(configPath string) (config *Configuration, err error) {
	if configPath == "" {
		configPath = defaultPath()
	}

	log.Info("Reading configuration file %s", configPath)
	configBuf, err := ioutil.ReadFile(configPath)
	if err != nil {
		err = fmt.Errorf("Failed to read configuration file %s: %v", configPath, err)
		return
	}
	config = new(Configuration)
	if err = yaml.Unmarshal(configBuf, &config); err != nil {
		err = fmt.Errorf("Error parsing configuration file %s: %v", configPath, err)
		return
	}

	if config.Domain == "" {
		config.Domain = defaultDomain
	}

	if config.LogTo == "" {
		config.LogTo = defaultLogto
	}

	if config.LogLevel == "" {
		config.LogLevel = defaultLogLevel
	}

	if config.TunnelAddr == "" {
		config.TunnelAddr = defaultTunnelAddr
	}

	return
}

func defaultPath() string {

	filename, _ := osext.Executable()
	dir := path.Dir(filename)

	return path.Join(dir, "ngrokd.yaml")
}
