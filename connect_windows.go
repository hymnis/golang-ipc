// go:build windows
package ipc

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/Microsoft/go-winio"
)

// Server function
// Create the named pipe (if it doesn't already exist) or tcp endpoint and start listening for a client to connect.
// when a client connects and connection is accepted the read function is called on a go routine.
func (sc *Server) run() error {

	var err error
	var listen net.Listener

	if sc.network {
		listen, err = net.Listen("tcp", fmt.Sprint(sc.networkListen)+":"+fmt.Sprint(sc.networkPort))
	} else {
		var pipeBase = `\\.\pipe\`

		listen, err = winio.ListenPipe(pipeBase+sc.name, nil)
	}

	if err != nil {
		return err
	}

	sc.listen = listen
	sc.status = Listening
	sc.recieved <- &Message{Status: sc.status.String(), MsgType: -1}
	sc.connChannel = make(chan bool)

	go sc.acceptLoop()

	err2 := sc.connectionTimer()
	if err2 != nil {
		return err2
	}

	return nil

}

// Client function
// dial - attempts to connect to a named pipe or tcp endpoint created by the server
func (cc *Client) dial() error {

	var net_type, address string
	var err error
	var pn net.Conn

	startTime := time.Now()

	for {
		if cc.timeout != 0 {
			if time.Now().Sub(startTime).Seconds() > cc.timeout {
				cc.status = Closed
				return errors.New("Timed out trying to connect")
			}
		}

		if cc.network {
			net_type = "tcp"
			address = fmt.Sprint(cc.networkServer) + ":" + fmt.Sprint(cc.networkPort)
			pn, err = net.Dial(net_type, address)
		} else {
			pipeBase := `\\.\pipe\`
			pn, err = winio.DialPipe(pipeBase+cc.Name, nil)
		}

		if err != nil {

			if strings.Contains(err.Error(), "The system cannot find the file specified.") {

			} else {
				return err
			}

		} else {

			cc.conn = pn
			err = cc.handshake()

			if err != nil {
				return err
			}

			return nil
		}

		time.Sleep(cc.retryTimer * time.Second)

	}
}
