package msg

import (
	"encoding/binary"
	"fmt"
	"mgrok/log"
	"net"
)

var logger = log.NewPrefixLogger("conn")

func readMsgShared(c net.Conn) (buffer []byte, err error) {
	logger.Debug("Waiting to read message")

	var sz int64
	err = binary.Read(c, binary.LittleEndian, &sz)
	if err != nil {
		return
	}
	logger.Debug("Reading message with length: %d", sz)

	buffer = make([]byte, sz)
	n, err := c.Read(buffer)
	logger.Debug("Read message %s", buffer)

	if err != nil {
		return
	}

	if int64(n) != sz {
		err = fmt.Errorf("Expected to read %d bytes, but only read %d", sz, n)
		// err = errors.New(fmt.Sprintf("Expected to read %d bytes, but only read %d", sz, n))
		return
	}

	return
}

//ReadMsg read message
func ReadMsg(c net.Conn) (msg Message, err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}

	return Unpack(buffer)
}

//ReadMsgInto read message info
func ReadMsgInto(c net.Conn, msg Message) (err error) {
	buffer, err := readMsgShared(c)
	if err != nil {
		return
	}
	return UnpackInto(buffer, msg)
}

//WriteMsg write message
func WriteMsg(c net.Conn, msg interface{}) (err error) {
	buffer, err := Pack(msg)
	if err != nil {
		return
	}

	logger.Debug("Writing message: %s", string(buffer))
	err = binary.Write(c, binary.LittleEndian, int64(len(buffer)))

	if err != nil {
		return
	}

	if _, err = c.Write(buffer); err != nil {
		return
	}

	return nil
}
