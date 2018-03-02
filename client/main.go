package client

import (
	"fmt"
	"math/rand"
	"mgrok/client/mvc"
	"mgrok/log"
	"mgrok/util"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/inconshreveable/mousetrap"
	"github.com/olekukonko/tablewriter"
)

func init() {
	if runtime.GOOS == "windows" {
		if mousetrap.StartedByExplorer() {
			fmt.Println("Don't double-click mgrok!")
			fmt.Println("You need to open cmd.exe and run it from the command line!")
			time.Sleep(5 * time.Second)
			os.Exit(1)
		}
	}
}

// Main clinet main function
func Main() {
	opts, err := ParseArgs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// set up logging
	log.LogTo(opts.log, opts.loglevel)

	// read configuration file
	config, err := LoadConfiguration(opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// seed random number generator
	seed, err := util.RandomSeed()
	if err != nil {
		fmt.Printf("Couldn't securely seed the random number generator!")
		os.Exit(1)
	}
	rand.Seed(seed)

	model := newClientModel(config)

	model.updateCallback = func(c *Model) {
		changedTime = time.Now()

		render(c)
		// fmt.Println("%s", time.Now())
	}

	model.Run()

}

var changedTime time.Time
var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func render(c *Model) {

	callClear()

	// size := len(c.tunnels) + 3
	// data := make([][]string, size)

	// i := 0
	// data[i] = make([]string, 2)
	// data[i][0] = "Tunnel Status"

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
	// i = i + 1

	// data[i] = make([]string, 2)
	// data[i][0] = "Version"
	// data[i][1] = fmt.Sprintf("%s/%s", c.GetClientVersion(), c.GetServerVersion())
	// i = i + 1

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

	fmt.Println()
	fmt.Println()
}
