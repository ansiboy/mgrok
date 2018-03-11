package client

import (
	"fmt"
	"mgrok/log"
	"time"

	"github.com/gdamore/tcell"

	"github.com/rivo/tview"
	// "github.com/gizak/termui"
)

var table *tview.Table
var app = tview.NewApplication()

func startConsole(modelChan chan *Model) error {
	go func() {
		for {
			c := <-modelChan
			render(c)
		}
	}()

	table = tview.NewTable().SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEscape || key == tcell.KeyCtrlC {
			app.Stop()
		}
	})

	cell := tview.NewTableCell("connecting").SetAlign(tview.AlignCenter)
	table.SetCell(0, 0, cell)
	app.SetRoot(table, true)
	err := app.Run()
	if err != nil {
		app = nil
		log.Error("init tview fail.\r")
	}
	return err
}

func render(c *Model) {
	var connStatus = ""
	switch c.connStatus {
	case ConnConnecting:
		connStatus = "connecting"
	case ConnReconnecting:
		connStatus = "connecting"
	case ConnOnline:
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
	for r := 0; r < len(data); r++ {
		table.SetCell(r, 0, tview.NewTableCell(data[r][0]))
		table.SetCell(r, 1, tview.NewTableCell(data[r][1]))
	}

	app.Draw()
}
