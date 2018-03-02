package client

import (
	"fmt"
	"mgrok/log"
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

const (
	fgColor = termbox.ColorWhite
	bgColor = termbox.ColorDefault
)

type area struct {
	// top-left corner
	x, y int

	// size of the area
	w, h int

	// default colors
	fgColor, bgColor termbox.Attribute
}

func NewArea(x, y, w, h int) *area {
	return &area{x, y, w, h, fgColor, bgColor}
}

func (a *area) Clear() {
	for i := 0; i < a.w; i++ {
		for j := 0; j < a.h; j++ {
			termbox.SetCell(a.x+i, a.y+j, ' ', a.fgColor, a.bgColor)
		}
	}
}

func (a *area) APrintf(fg termbox.Attribute, x, y int, arg0 string, args ...interface{}) {
	s := fmt.Sprintf(arg0, args...)
	for i, ch := range s {
		termbox.SetCell(a.x+x+i, a.y+y, ch, fg, bgColor)
	}
}

func (a *area) Printf(x, y int, arg0 string, args ...interface{}) {
	a.APrintf(a.fgColor, x, y, arg0, args...)
}

const (
	size          = 10
	pathMaxLength = 25
)

type HttpView struct {
	log.Logger
	*area

	// httpProto    *proto.Http
	// HttpRequests *util.Ring
	shutdown chan int
	// termView     *TermView
}

func colorFor(status string) termbox.Attribute {
	switch status[0] {
	case '3':
		return termbox.ColorCyan
	case '4':
		return termbox.ColorYellow
	case '5':
		return termbox.ColorRed
	default:
	}
	return termbox.ColorWhite
}

func newTermHttpView(x, y int) *HttpView { // termView *TermView, proto *proto.Http,
	v := &HttpView{
		// httpProto:    proto,
		// HttpRequests: util.NewRing(size),
		area:     NewArea(x, y, 70, size+5),
		shutdown: make(chan int),
		// termView:     termView,
		Logger: log.NewPrefixLogger("view", "term", "http"),
	}
	// ctl.Go(v.Run)
	return v
}

func (v *HttpView) Run() {
	// updates := v.httpProto.Txns.Reg()

	for {
		// select {
		// case txn := <-updates:
		// v.Debug("Got HTTP update")
		// if txn.(*proto.HttpTxn).Resp == nil {
		// v.HttpRequests.Add(txn)
		// }
		// v.Render()
		//  }
	}

}

func (v *HttpView) Render() {
	v.Clear()
	v.Printf(0, 0, "HTTP Requests")
	v.Printf(0, 1, "-------------")
	// for i, obj := range v.HttpRequests.Slice() {
	i := 0
	// txn := "http" //obj.(*proto.HttpTxn)
	// path := ""    //truncatePath(txn.Req.URL.Path)
	v.Printf(0, 3+i, "%s %v", "txn.Req.Method", "path")
	// if txn.Resp != nil {
	// 	v.APrintf(colorFor(txn.Resp.Status), 30, 3+i, "%s", txn.Resp.Status)
	// }
	// }
	// termbox.Flush()
}

func (v *HttpView) Shutdown() {
	close(v.shutdown)
}

func truncatePath(path string) string {
	// Truncate all long strings based on rune count
	if utf8.RuneCountInString(path) > pathMaxLength {
		path = string([]rune(path)[:pathMaxLength])
	}

	// By this point, len(path) should be < pathMaxLength if we're dealing with single-byte runes.
	// Otherwise, we have a multi-byte string and need to calculate the size of each rune and
	// truncate manually.
	//
	// This is a workaround for a bug in termbox-go. Remove it when this issue is fixed:
	// https://github.com/nsf/termbox-go/pull/21
	if len(path) > pathMaxLength {
		out := make([]byte, pathMaxLength, pathMaxLength)
		length := 0
		for {
			r, size := utf8.DecodeRuneInString(path[length:])
			if r == utf8.RuneError && size == 1 {
				break
			}

			// utf8.EncodeRune expects there to be enough room to store the full size of the rune
			if length+size <= pathMaxLength {
				utf8.EncodeRune(out[length:], r)
				length += size
			} else {
				break
			}
		}
		path = string(out[:length])
	}
	return path
}
