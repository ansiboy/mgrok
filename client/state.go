package client

type UpdateStatus int

const (
	UpdateNone = -1 * iota
	UpdateInstalling
	UpdateReady
	UpdateAvailable
)

type ConnStatus int

const (
	ConnConnecting = iota
	ConnReconnecting
	ConnOnline
)

type Tunnel struct {
	PublicUrl string
	LocalAddr string
	Type      string
}
