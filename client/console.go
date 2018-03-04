package client

import (
	"fmt"
	"mgrok/client/mvc"
	"mgrok/log"
	"time"

	"github.com/gizak/termui"
)

func startConsole(modelChan chan *Model) {
	err := termui.Init()
	if err != nil {
		// panic(err)
		log.Error("init termui fail.\r")
	}
	defer termui.Close()
	go func() {
		for {
			c := <-modelChan
			render(c)
		}
	}()

	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		// handle Ctrl + x combination
		termui.StopLoop()
	})

	termui.Loop()
}

func render(c *Model) {
	var connStatus = ""
	switch c.connStatus {
	case mvc.ConnConnecting:
		connStatus = "connecting"
	case mvc.ConnReconnecting:
		connStatus = "connecting"
	case mvc.ConnOnline:
		connStatus = "online"
	}

	var version = fmt.Sprintf("%s/%s", c.GetClientVersion(), c.GetServerVersion())

	connCount := c.metrics.connMeter.Count()

	msec := float64(time.Millisecond)
	avgConnTime := float64(c.metrics.connTimer.Mean() / msec)

	data := [][]string{
		[]string{"mgrok", "Ctrl+C to quit"},
		[]string{"Tunnel Status", connStatus},
		[]string{"Version", version},
		[]string{"# Conn", fmt.Sprintf("%d", connCount)},
		[]string{"Avg Conn Time", fmt.Sprintf("%.2fms", avgConnTime)},
	}

	size := len(c.tunnels)
	tunnels := make([][]string, size)

	i := 0
	for _, t := range c.tunnels {
		tunnels[i] = make([]string, 2)
		tunnels[i][0] = "Forwarding"
		tunnels[i][1] = fmt.Sprintf("%s -> %s", t.PublicUrl, t.LocalAddr)
		i = i + 1
	}

	data = append(data, tunnels...)

	table1 := termui.NewTable()
	table1.Rows = data
	table1.FgColor = termui.ColorWhite
	table1.BgColor = termui.ColorDefault
	table1.Y = 0
	table1.X = 0

	table1.Analysis()
	table1.SetSize()

	termui.Render(table1)
}
