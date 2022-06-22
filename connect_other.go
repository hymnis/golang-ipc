//go:build linux || darwin
// +build linux darwin

package ipc

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"syscall"
	"time"
)

// Server create a unix socket or tcp endpoint and start listening for connections - for unix and linux
func (sc *Server) run() error {
	var oldUmask int
	var net_type, address string

	if sc.network {
		net_type = "tcp"
		address = "0.0.0.0:" + fmt.Sprint(sc.networkPort)
	} else {
		base := "/tmp/"
		sock := ".sock"

		if err := os.RemoveAll(base + sc.name + sock); err != nil {
			return err
		}

		if sc.unMask {
			oldUmask = syscall.Umask(0)
		}

		net_type = "unix"
		address = base + sc.name + sock
	}

	listen, err := net.Listen(net_type, address)

	if !sc.network && sc.unMask {
		syscall.Umask(oldUmask)
	}

	if err != nil {
		return err
	}

	sc.listen = listen
	sc.status = Listening
	sc.recieved <- &Message{Status: sc.status.String(), MsgType: -1}
	sc.connChannel = make(chan bool)

	go sc.acceptLoop()

	err = sc.connectionTimer()
	if err != nil {
		return err
	}

	return nil

}

// Client connect to the unix socket or tcp endpoint created by the server - for unix and linux
func (cc *Client) dial() error {

	var net_type, address string

	startTime := time.Now()

	for {
		if cc.timeout != 0 {
			if time.Since(startTime).Seconds() > cc.timeout {
				cc.status = Closed
				return errors.New("Timed out trying to connect")
			}
		}

		if cc.network {
			net_type = "tcp"
			address = "0.0.0.0:" + fmt.Sprint(cc.networkPort)
		} else {
			base := "/tmp/"
			sock := ".sock"
			net_type = "unix"
			address = base + cc.Name + sock
		}

		conn, err := net.Dial(net_type, address)
		if err != nil {

			if strings.Contains(err.Error(), "connect: no such file or directory") {

			} else if strings.Contains(err.Error(), "connect: connection refused") {

			} else {
				cc.recieved <- &Message{err: err, MsgType: -2}
			}

		} else {

			cc.conn = conn

			err = cc.handshake()
			if err != nil {
				return err
			}

			return nil
		}

		time.Sleep(cc.retryTimer * time.Second)

	}

}
