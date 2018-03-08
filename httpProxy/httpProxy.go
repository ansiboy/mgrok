package httpProxy

import (
	"encoding/json"
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
var logger log.Logger // =
func Main() {
	opts := parseArgs()
	config, err := loadConfiguration(opts.config)
	checkError(err)
	if config.LogTo != "" {
		log.LogTo(config.LogTo, config.LogLevel)
		logger = log.NewPrefixLogger("httpProxy")
	}

	go start(config)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		client := new(http.Client)

		var err error

		redirectInfo, ok := tunnleInfos[request.Host]
		if ok == false {
			http.NotFound(writer, request)
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

		for key, values := range respons.Header {
			for _, value := range values {
				writer.Header().Set(key, value)
			}
		}

		body, _ := ioutil.ReadAll(respons.Body)
		writer.Write(body)

	})

	logger.Info("HTTP service listen at %s", config.HTTPAddr)
	http.ListenAndServe(config.HTTPAddr, nil)
}

func internalServerError(writer http.ResponseWriter) {
	msg := string(http.StatusInternalServerError) + " internal server error"
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
			logger.Info("Delete tunnel registry")
			delete(tunnleInfos, info.SourceAddr)
		case ActionRegister:
			logger.Info("Register tunnel")
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
