package httpProxy

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
)

func Main() {

	config, err := loadConfiguration("")
	checkError(err)

	go startData(config)

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		client := new(http.Client)

		var err error

		targetURL := "http://127.0.0.1:8080" + request.URL.Path
		request.URL, err = url.Parse(targetURL)
		checkError(err)

		request.RequestURI = ""

		respons, err := client.Do(request)
		checkError(err)

		body, _ := ioutil.ReadAll(respons.Body)

		writer.Header()
		writer.Write(body)

	})

	http.ListenAndServe(":3762", nil)

}

func startData(config *Configuration) {
	listener, err := net.Listen("tcp", config.dataAddr)
	checkError(err)
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		checkError(err)

		for {
			data, err := ioutil.ReadAll(conn)
			if err != nil || len(data) == 0 {
				continue
			}

			fmt.Println(string(data))
		}

	}
}

func checkError(err error) {
	if err == nil {
		return
	}

	panic(err)
}

func handleConn(conn net.Conn) {

	target, err := net.Dial("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	const bufferSize = 1024
	buf := make([]byte, bufferSize)
	for {
		length, err := conn.Read(buf)
		if err != nil {
			panic(err)
		}

		target.Write(buf)
		if length < bufferSize {
			break
		}
	}
}
