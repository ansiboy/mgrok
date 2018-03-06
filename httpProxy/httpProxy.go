package httpProxy

import (
	"encoding/json"
	"io/ioutil"
	"mgrok/log"
	"net"
	"net/http"
	"net/url"
)

type HTTPRedirect struct {
	SourceAddr string
	TargetAddr string
}

type ActionData struct {
	Action string
	Data   HTTPRedirect
}

var ACTION_PING = "PING"
var ACTION_PONG = "PONG"
var ACTION_REGISTER = "REGISTER"
var ACTION_DELETE = "DELETE"
var tunnleInfos = make(map[string]HTTPRedirect)
var logger log.Logger // =
func Main() {

	config, err := loadConfiguration("")
	checkError(err)
	if config.LogTo != "" {
		log.LogTo(config.LogTo, config.LogLevel)
		logger = log.NewPrefixLogger("httpProxy")
	}

	go start(config)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		client := new(http.Client)

		var err error

		redirectInfo, ok := tunnleInfos[request.RemoteAddr]
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
	listener, err := net.Listen("tcp", config.DataAddr)
	logger.Info("Data service listen at %s", config.DataAddr)

	checkError(err)
	defer func() {
		log.Info("Data service listener close")
		listener.Close()
	}()

	conns := make([]net.Conn, 100)
	defer func() {
		for _, conn := range conns {
			if conn != nil {
				conn.Close()
			}
		}
	}()

	for {
		conn, err := listener.Accept()
		logger.Info("Connected for client %", conn.RemoteAddr())
		checkError(err)
		_ = append(conns, conn)

		for {
			data, err := ioutil.ReadAll(conn)
			if err != nil || len(data) == 0 {
				continue
			}

			var f ActionData
			err = json.Unmarshal(data, f)
			if err != nil {
				continue
			}

			action := f.Action
			info := f.Data
			switch action {
			case ACTION_DELETE:
				logger.Info("Delete tunnel registry")
				delete(tunnleInfos, info.SourceAddr)
			case ACTION_REGISTER:
				logger.Info("Register tunnel")
				tunnleInfos[info.SourceAddr] = info
			}
		}
	}
	// }()
}

func checkError(err error) {
	if err == nil {
		return
	}

	panic(err)
}
