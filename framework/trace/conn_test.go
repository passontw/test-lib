// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package trace

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// ConnTCPListener takes a TCP listener and accepts n TCP connections
// Returns connections using connChan
func connTCPListener(t *testing.T, n int, ln net.Listener, connChan chan<- net.Conn) {
	// Listen and accept n incoming connections
	for i := 0; i < n; i++ {
		conn, err := ln.Accept()
		if err != nil {
			t.Log("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		// Send accepted connection to channel
		connChan <- conn
	}
	ln.Close()
	close(connChan)
}

func TestConn(t *testing.T) {
	log := NewLogger(1000)
	log.SetLogger("conn", `{"network":"tcp","addr":":7020"}`)
	log.Informational("informational")
}

// need to rewrite this test, it's not stable
func TestReconnect(t *testing.T) {
	// Setup connection listener
	newConns := make(chan net.Conn)
	connNum := 2
	ln, err := net.Listen("tcp", ":6002")
	if err != nil {
		t.Log("Error listening:", err.Error())
		os.Exit(1)
	}
	go connTCPListener(t, connNum, ln, newConns)

	// Setup trace
	log := NewLogger(1000)
	log.SetPrefix("test")
	log.SetLogger(AdapterConn, `{"network":"tcp","reconnect":true,"level":6,"addr":":6002"}`)
	log.Informational("informational 1")

	// Refuse first connection
	first := <-newConns
	first.Close()

	// Send another trace after conn closed
	log.Informational("informational 2")

	// Check if there was a second connection attempt
	select {
	case second := <-newConns:
		second.Close()
	default:
		t.Error("Did not reconnect")
	}
}

func TestConnWriter_Format(t *testing.T) {
	lg := &LogMsg{
		Level:      LevelDebug,
		Msg:        "Hello, world",
		When:       time.Date(2020, 9, 19, 20, 12, 37, 9, time.UTC),
		FilePath:   "/user/home/main.go",
		LineNumber: 13,
		Prefix:     "Cus",
	}
	cw := NewConn().(*connWriter)
	res := cw.Format(lg)
	assert.Equal(t, "[D] Cus Hello, world", res)
}
