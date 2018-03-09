package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mgrok/log"
	"mgrok/proxy"
	"net/http"
)

// TunnelCenterRegistry 遂道注册中心，负责将遂道信息添加到中心的代理服务器
type TunnelCenterRegistry struct {
	// conn       net.Conn
	url        string
	targetAddr string
}

var logger = log.NewPrefixLogger("RedirectData")

//var url = "http://" + dataAddr + "/control"

func newRedirectData(dataAddr string, targetAddr string) (registry *TunnelCenterRegistry, err error) {
	// var conn net.Conn
	// conn, err = http.Client.Post(url) //http.Dial(url, "", "http://"+dataAddr+"/ws")
	// if err != nil {
	// 	return
	// }

	registry = &TunnelCenterRegistry{
		url:        "http://" + dataAddr + "/control",
		targetAddr: targetAddr,
	}

	// go func() {

	// 	// defer conn.Close()

	// 	for {
	// 		// data, err := ioutil.ReadAll(conn)
	// 		if err != nil {
	// 			//TODO 处理异常
	// 			continue
	// 		}

	// 		var f httpProxy.ActionData
	// 		err = json.Unmarshal(data, f)
	// 		if err != nil {
	// 			//TODO 处理异常
	// 			continue
	// 		}

	// 		switch f.Action {

	// 		}

	// 		fmt.Println(data)
	// 	}
	// }()

	return
}

func (r *TunnelCenterRegistry) register(url string, t *Tunnel) {
	//当前仅处理 http 的转发
	addr := cutHTTPPrefix(url)
	if addr == "" {
		return
	}

	actionData := httpProxy.ActionData{
		Action: httpProxy.ActionRegister,
		Data: httpProxy.HTTPRedirect{
			SourceAddr: addr,
			TargetAddr: r.targetAddr,
		},
	}

	r.write(actionData)
	// data, err := json.Marshal(actionData)
	// if err != nil {
	// 	logger.Error(err.Error(), err)
	// 	return
	// }

	// _, err = r.conn.Write(data)
	// if err != nil {
	// 	logger.Error(err.Error())
	// }
}

func (r *TunnelCenterRegistry) del(url string) {
	addr := cutHTTPPrefix(url)
	if addr == "" {
		return
	}
	actionData := httpProxy.ActionData{
		Action: httpProxy.ActionDelete,
		Data: httpProxy.HTTPRedirect{
			SourceAddr: addr,
			TargetAddr: r.targetAddr,
		},
	}

	r.write(actionData)
}

func (r *TunnelCenterRegistry) write(actionData httpProxy.ActionData) {
	data, err := json.Marshal(actionData)
	fmt.Println(string(data))
	if err != nil {
		logger.Error(err.Error(), err)
		return
	}

	_, err = http.Post(r.url, "text/application-json", bytes.NewBuffer(data))
	if err != nil {
		logger.Error(err.Error())
	}
}

func cutHTTPPrefix(url string) string {
	const prefix = "http://"
	prefixLength := len(prefix)
	if url[0:prefixLength] != prefix {
		return ""
	}

	return url[prefixLength:]
}
