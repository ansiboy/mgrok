package client

import (
	"fmt"
	"math/rand"
	"mgrok/log"
	"mgrok/util"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/inconshreveable/mousetrap"
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

	modelChan := make(chan *Model)
	defer close(modelChan)

	model := newClientModel(config, modelChan)

	go model.Run()

	if config.PprofAddr != "" {
		go func() {
			http.ListenAndServe(config.PprofAddr, nil)
		}()
	}

	if config.LogTo != "stdout" {
		startConsole(model.changed)
		return
	}

	for {
		c := <-model.changed
		// 延时 0.5 秒，让其他信息先输出
		// time.AfterFunc(500, func() {
		// fmt.Println()
		printModelInfo(c)
		// fmt.Println()
		// })
	}
}
