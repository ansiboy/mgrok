package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mgrok/httpProxy"
	"mgrok/log"
	"net"
)

// TunnelCenterRegistry 遂道注册中心，负责将遂道信息添加到中心的代理服务器
type TunnelCenterRegistry struct {
	conn       net.Conn
	targetAddr string
}

func newRedirectData(dataAddr string, targetAddr string) (registry *TunnelCenterRegistry, err error) {
	var conn net.Conn
	conn, err = net.Dial("tcp", dataAddr)
	if err != nil {
		return
	}

	registry = &TunnelCenterRegistry{
		conn:       conn,
		targetAddr: targetAddr,
	}

	go func() {

		defer conn.Close()

		for {
			data, err := ioutil.ReadAll(conn)
			if err != nil {
				//TODO 处理异常
				continue
			}

			var f httpProxy.ActionData
			err = json.Unmarshal(data, f)
			if err != nil {
				//TODO 处理异常
				continue
			}

			switch f.Action {

			}

			fmt.Println(data)
		}
	}()

	return
}

func (r *TunnelCenterRegistry) register(url string, t *Tunnel) {
	//当前仅处理 http 的转发
	addr := cutHTTPPrefix(url)
	if addr == "" {
		return
	}

	actionData := httpProxy.ActionData{
		Action: httpProxy.ACTION_REGISTER,
		Data: httpProxy.HTTPRedirect{
			SourceAddr: addr,
			TargetAddr: r.targetAddr,
		},
	}
	data, err := json.Marshal(actionData)
	if err != nil {
		log.Error(err.Error(), err)
		return
	}

	r.conn.Write(data)
}

func (r *TunnelCenterRegistry) del(url string) {
	addr := cutHTTPPrefix(url)
	if addr == "" {
		return
	}
	actionData := httpProxy.ActionData{
		Action: httpProxy.ACTION_DELETE,
		Data: httpProxy.HTTPRedirect{
			SourceAddr: addr,
			TargetAddr: r.targetAddr,
		},
	}

	data, err := json.Marshal(actionData)
	if err != nil {
		log.Error(err.Error(), err)
		return
	}

	r.conn.Write(data)
}

func cutHTTPPrefix(url string) string {
	const prefix = "http://"
	prefixLength := len(prefix)
	if url[0:prefixLength] != prefix {
		return ""
	}

	return url[prefixLength:]
}
