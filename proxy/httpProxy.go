package httpProxy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mgrok/log"
	"net/http"
	"net/url"
)

type HTTPRedirect struct {
	SourceAddr string `json:"sourceAddr"`
	TargetAddr string `json:"targetAddr"`
}

type ActionData struct {
	Action string       `json:"action"`
	Data   HTTPRedirect `json:"data"`
}

// ActionRegister 注册隧道
var ActionRegister = "REGISTER"

// ActionDelete 删除隧道
var ActionDelete = "DELETE"

var tunnleInfos = make(map[string]HTTPRedirect)
var logger log.Logger

const statusCodeTunnelNotFound = 550

func Main() {
	opts := parseArgs()
	config, err := loadConfiguration(opts.config)
	checkError(err)
	log.LogTo(config.LogTo, config.LogLevel)
	logger = log.NewPrefixLogger("mgrokp")
	go start(config)

	server := &http.Server{
		Addr: config.HTTPAddr,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			client := new(http.Client)

			var err error

			host := request.Header.Get("X-Host")
			if len(host) == 0 {
				host = request.Host
			}
			redirectInfo, ok := tunnleInfos[host]
			if ok == false {
				msg := fmt.Sprintf("Tunnel %s not found", host)
				http.Error(writer, msg, statusCodeTunnelNotFound)
				// fmt.Fprintf(writer, "Tunnel %s not found", host)
				// fmt.Printf("Tunnel %s not found", host)
				return
			}

			targetURL := "http://" + redirectInfo.TargetAddr + request.URL.Path
			request.URL, err = url.Parse(targetURL)
			if err != nil {
				log.Error(err.Error())
				internalServerError(writer)
				return
			}

			request.RequestURI = ""
			respons, err := client.Do(request)
			if err != nil {
				log.Error(err.Error())
				internalServerError(writer)
				return
			}

			defer respons.Body.Close()

			writer.WriteHeader(respons.StatusCode)
			for key, values := range respons.Header {
				for _, value := range values {
					writer.Header().Set(key, value)
				}
			}

			body, _ := ioutil.ReadAll(respons.Body)
			writer.Write(body)
		}),
	}

	if config.PprofAddr != "" {
		go func() {
			http.ListenAndServe(config.PprofAddr, nil)
		}()
	}

	logger.Info("HTTP service listen at %s", config.HTTPAddr)
	server.ListenAndServe()
}

func internalServerError(writer http.ResponseWriter) {
	msg := fmt.Sprintf("%d %s", http.StatusInternalServerError, "internal server error") // string(http.StatusInternalServerError) +
	http.Error(writer, msg, http.StatusInternalServerError)
}

func start(config *Configuration) {

	logger.Info("Data service listen at %s", config.DataAddr)
	http.HandleFunc("/control", func(writer http.ResponseWriter, request *http.Request) {

		data, err := ioutil.ReadAll(request.Body)

		if err != nil {
			logger.Error("Read data from connection fail", err)
			return
		}

		obj := &ActionData{}
		err = json.Unmarshal(data, obj)
		if err != nil {
			logger.Error("Parse json object fail")
			logger.Error(string(data))
			return
		}

		action := obj.Action
		info := obj.Data
		switch action {
		case ActionDelete:
			logger.Info(fmt.Sprint("Delete tunnel registry ", info.SourceAddr))
			delete(tunnleInfos, info.SourceAddr)
		case ActionRegister:
			logger.Info(fmt.Sprint("Register tunnel ", info.SourceAddr))
			tunnleInfos[info.SourceAddr] = info
		}
	})
	http.ListenAndServe(config.DataAddr, nil)
}

func checkError(err error) {
	if err == nil {
		return
	}

	panic(err)
}
