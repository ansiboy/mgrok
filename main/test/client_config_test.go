package ngrok_test

import (
	"fmt"
	"mgrok/client"
	_ "mgrok/client"
	"testing"
)

func Test_LoadConfigration(t *testing.T) {
	opts, _ := client.ParseArgs()
	config, _ := client.LoadConfiguration(opts)
	if config == nil {
		t.Error("config is nil")
	}
	fmt.Printf("server_addr:%s\r\n", config.ServerAddr)
	fmt.Printf("trust_host_root_certs:%s\r\n", fmt.Sprint(config.TrustHostRootCerts))
}
