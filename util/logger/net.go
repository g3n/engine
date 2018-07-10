// Copyright 2016 The G3N Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logger

import (
	"net"
)

// Net is a network writer used for logging.
type Net struct {
	conn net.Conn
}

// NewNet creates and returns a pointer to a new Net object along with any error that occurred.
func NewNet(network string, address string) (*Net, error) {

	n := new(Net)
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	n.conn = conn
	return n, nil
}

// Write writes the provided logger event to the network.
func (n *Net) Write(event *Event) {

	n.conn.Write([]byte(event.fmsg))
}

// Clone closes the network connection.
func (n *Net) Close() {

	n.conn.Close()
}

func (n *Net) Sync() {

}
