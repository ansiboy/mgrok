package client

import (
	"ngrok/client/mvc"
	"ngrok/log"
)

type command interface{}

type cmdQuit struct {
	// display this message after quit
	message string
}

type cmdPlayRequest struct {
	// the tunnel to play this request over
	tunnel mvc.Tunnel

	// the bytes of the request to issue
	payload []byte
}

// The MVC Controller
type Controller struct {
	// Controller logger
	log.Logger

	// the model sends updates through this broadcast channel
	// updates *util.Broadcast

	// the model
	model *ClientModel

	// the views
	// views []mvc.View

	// internal structure to issue commands to the controller
	// cmds chan command

	// internal structure to synchronize access to State object
	// state chan *ClientModel

	// options
	// config *Configuration
}

// public interface
func NewController() *Controller {
	ctl := &Controller{
		Logger: log.NewPrefixLogger("controller"),
		// updates: util.NewBroadcast(),
		// cmds:    make(chan command),
		// views:   make([]mvc.View, 0),
		// state: make(chan *ClientModel),
	}

	return ctl
}

// func (ctl *Controller) State() *ClientModel {
// 	return <-ctl.state
// }

// func (ctl *Controller) Update(state *ClientModel) {
// 	ctl.updates.In() <- state
// }

// func (ctl *Controller) Updates() *util.Broadcast {
// 	return ctl.updates
// }

// func (ctl *Controller) Shutdown(message string) {
// 	ctl.cmds <- cmdQuit{message: message}
// }

// func (ctl *Controller) PlayRequest(tunnel mvc.Tunnel, payload []byte) {
// 	ctl.cmds <- cmdPlayRequest{tunnel: tunnel, payload: payload}
// }

// func (ctl *Controller) Go(fn func()) {
// 	go func() {
// 		defer func() {
// 			if r := recover(); r != nil {
// 				err := util.MakePanicTrace(r)
// 				ctl.Error(err)
// 				ctl.Shutdown(err)
// 			}
// 		}()

// 		fn()
// 	}()
// }

// private functions
/*
func (ctl *Controller) doShutdown() {
	ctl.Info("Shutting down")

	var wg sync.WaitGroup

	// wait for all of the views, plus the model
	// wg.Add(len(ctl.views) + 1)

	// for _, v := range ctl.views {
	// 	vClosure := v
	// 	ctl.Go(func() {
	// 		vClosure.Shutdown()
	// 		wg.Done()
	// 	})
	// }

	// ctl.Go(func() {
	// 	ctl.model.Shutdown()
	// 	wg.Done()
	// })

	go func() {
		ctl.model.Shutdown()
		wg.Done()
	}()

	wg.Wait()
}
*/

// func (ctl *Controller) AddView(v mvc.View) {
// 	ctl.views = append(ctl.views, v)
// }

// func (ctl *Controller) GetWebInspectAddr() string {
// 	return ctl.config.InspectAddr
// }

func (ctl *Controller) SetupModel(config *Configuration) *ClientModel {
	model := newClientModel(config, ctl)
	ctl.model = model
	return model
}

// func (ctl *Controller) GetModel() *ClientModel {
// 	return ctl.model //ctl.model.(*ClientModel)
// }

func (ctl *Controller) Run(config *Configuration) {
	// Save the configuration

	// ctl.config = config
	ctl.SetupModel(config)

	done := make(chan *ClientModel)

	go ctl.model.Run()

	// ctl.state <- ctl.model
	done <- ctl.model
}
