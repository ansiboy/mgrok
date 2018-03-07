package httpProxy

import (
	"fmt"
	"io/ioutil"
	"mgrok/log"
	"path"

	"github.com/kardianos/osext"
	yaml "gopkg.in/yaml.v1"
)

const (
	defaultHTTPAddr = "127.0.0.1:3762"
	defaultDataAddr = "127.0.0.1:6523"
)

// Configuration http proxy configuration
type Configuration struct {
	HTTPAddr string `yaml:"http_addr,omitempty"`
	DataAddr string `yaml:"data_addr,omitempty"`
	LogTo    string `yaml:"log_to,omitempty"`
	LogLevel string `yaml:"log_level,omitempty"`
}

func loadConfiguration(configPath string) (config *Configuration, err error) {
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
	if err = yaml.Unmarshal(configBuf, config); err != nil {
		err = fmt.Errorf("Error parsing configuration file %s: %v", configPath, err)
		return
	}

	if config.DataAddr == "" {
		config.DataAddr = defaultDataAddr
	}

	if config.HTTPAddr == "" {
		config.HTTPAddr = defaultHTTPAddr
	}

	return
}

func defaultPath() string {

	filename, _ := osext.Executable()
	dir := path.Dir(filename)

	return path.Join(dir, "httpProxy.yaml")
}
