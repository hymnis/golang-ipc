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
func (sc *Server) runSocket() error {
	var pipeBase = `\\.\pipe\`

	listen, err := winio.ListenPipe(pipeBase+sc.name, nil)

	if err != nil {

		return err
	}

	return run(listen)
}

func (sc *Server) runNetwork() error {
	listen, err := net.Listen("tcp", "0.0.0.0:"+fmt.Sprint(sc.networkPort))

	if err != nil {

		return err
	}

	return run(listen)
}

// Create the named pipe (if it doesn't already exist) or tcp endpoint and start listening for a client to connect.
// when a client connects and connection is accepted the read function is called on a go routine.
func run(sc *Server, listen *Listener) error {

	sc.listen = listen

	sc.status = Listening
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
			address = "0.0.0.0:" + fmt.Sprint(cc.networkPort)
			pn, err := net.Dial(net_type, address)
		} else {
			pipeBase := `\\.\pipe\`
			pn, err := winio.DialPipe(pipeBase+cc.Name, nil)
		}

		if err != nil {

			if strings.Contains(err.Error(), "The system cannot find the file specified.") == true {

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
