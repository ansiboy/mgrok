package proto

type Protocol interface {
	GetName() string
	// WrapConn(net.Conn, interface{}) net.Conn
}
