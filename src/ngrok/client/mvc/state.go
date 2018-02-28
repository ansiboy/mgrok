package mvc

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
	// Protocol  proto.Protocol
	LocalAddr string
	Type      string
}

type ConnectionContext struct {
	Tunnel     Tunnel
	ClientAddr string
}

// type State interface {
// 	GetClientVersion() string
// 	GetServerVersion() string
// 	GetTunnels() []Tunnel
// 	GetProtocols() []proto.Protocol
// 	GetUpdateStatus() UpdateStatus
// 	GetConnStatus() ConnStatus
// 	GetConnectionMetrics() (metrics.Meter, metrics.Timer)
// 	GetBytesInMetrics() (metrics.Counter, metrics.Histogram)
// 	GetBytesOutMetrics() (metrics.Counter, metrics.Histogram)
// 	SetUpdateStatus(UpdateStatus)
// }
