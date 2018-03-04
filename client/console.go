package client

import (
	"fmt"
	"mgrok/client/mvc"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/olekukonko/tablewriter"
)

//create a map for storing clear funcs

type Console struct {
	clear map[string]func()
	model *Model
}

func NewConsole(model *Model) Console {
	c := Console{
		clear: make(map[string]func()),
		model: model,
	}

	c.clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c.clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	c.clear["darwin"] = c.clear["linux"]

	return c
}

func (c Console) callClear() {
	goos := runtime.GOOS
	value, ok := c.clear[goos] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                    //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic(fmt.Sprintf("Your platform %s is unsupported! I can't clear terminal screen :(", goos))
	}
}

func (console Console) Render() {
	c := console.model
	console.callClear()

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

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"mgrok", "Ctrl+C to quit"})
	table.SetColWidth(200)
	table.SetColWidth(600)
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
	table.Render()
}

// var clear map[string]func()

// func init() {
// 	clear = make(map[string]func()) //Initialize it
// 	clear["linux"] = func() {
// 		cmd := exec.Command("clear") //Linux example, its tested
// 		cmd.Stdout = os.Stdout
// 		cmd.Run()
// 	}
// 	clear["windows"] = func() {
// 		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
// 		cmd.Stdout = os.Stdout
// 		cmd.Run()
// 	}
// }

// func callClear() {
// 	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
// 	if ok {                          //if we defined a clear func for that platform:
// 		value() //we execute it
// 	} else { //unsupported platform
// 		panic("Your platform is unsupported! I can't clear terminal screen :(")
// 	}
// }

// func render(c *Model) {

// 	callClear()

// 	var connStatus = ""
// 	switch c.connStatus {
// 	case mvc.ConnConnecting:
// 		connStatus = "connecting"
// 	case mvc.ConnReconnecting:
// 		connStatus = "connecting"
// 	case mvc.ConnOnline:
// 		connStatus = "online"
// 	}

// 	var version = fmt.Sprintf("%s/%s", c.GetClientVersion(), c.GetServerVersion())

// 	connCount := c.metrics.connMeter.Count()

// 	msec := float64(time.Millisecond)
// 	avgConnTime := float64(c.metrics.connTimer.Mean() / msec)

// 	data := [][]string{
// 		[]string{"Tunnel Status", connStatus},
// 		[]string{"Version", version},
// 		[]string{"# Conn", fmt.Sprintf("%d", connCount)},
// 		[]string{"Avg Conn Time", fmt.Sprintf("%.2fms", avgConnTime)},
// 	}

// 	size := len(c.tunnels)
// 	tunnels := make([][]string, size)

// 	i := 0
// 	for _, t := range c.tunnels {
// 		tunnels[i] = make([]string, 2)
// 		tunnels[i][0] = "Forwarding"
// 		tunnels[i][1] = fmt.Sprintf("%s -> %s", t.PublicUrl, t.LocalAddr)
// 		i = i + 1
// 	}

// 	data = append(data, tunnels...)

// 	table := tablewriter.NewWriter(os.Stdout)
// 	table.SetHeader([]string{"mgrok", "Ctrl+C to quit"})
// 	table.SetColWidth(200)
// 	table.SetColWidth(600)
// 	table.SetRowLine(true)
// 	table.AppendBulk(data)
// 	table.SetColumnAlignment([]int{tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})
// 	table.Render()
// }
