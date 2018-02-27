package ngrokd_test

import (
	"fmt"
	// "ngrok/server"
	server "ngrok/server"
	"testing"
)

func Test_LoadConfigration(t *testing.T) {
	fmt.Println("Begin test")
	configPath := "/Volumes/data/projects/mgrok/src/ngrok/main/ngrok/ngrok.yaml"
	config, err := server.LoadConfiguration(configPath)
	if err != nil {
		fmt.Print(err)
		t.Error(err)
		return
	}
	if config == nil {
		t.Error("config is nil")
		return
	}
	fmt.Printf("http_addr %s\r\n", config.HttpAddr)
	fmt.Printf("https_addr %s\r\n", config.HttpsAddr)
	fmt.Printf("domain %s\r\n", config.Domain)
	fmt.Println("End test")
}
