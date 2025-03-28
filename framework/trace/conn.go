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
	"encoding/json"
	"fmt"
	"io"
	"net"

	"github.com/pkg/errors"
)

// connWriter implements LoggerInterface.
// Writes messages in keep-live tcp connection.
type connWriter struct {
	lg             *logWriter
	innerWriter    io.WriteCloser
	formatter      LogFormatter
	Formatter      string `json:"formatter"`
	ReconnectOnMsg bool   `json:"reconnectOnMsg"`
	Reconnect      bool   `json:"reconnect"`
	Net            string `json:"network"`
	Addr           string `json:"addr"`
	Level          int    `json:"level"`
}

// NewConn creates new ConnWrite returning as LoggerInterface.
func NewConn() Logger {
	conn := new(connWriter)
	conn.Level = LevelTrace
	conn.formatter = conn
	return conn
}

func (c *connWriter) Format(lm *LogMsg) string {
	return lm.OldStyleFormat()
}

// Init initializes a connection writer with json config.
// json config only needs they "level" key
func (c *connWriter) Init(config string) error {
	res := json.Unmarshal([]byte(config), c)
	if res == nil && len(c.Formatter) > 0 {
		fmtr, ok := GetFormatter(c.Formatter)
		if !ok {
			return errors.New(fmt.Sprintf("the formatter with name: %s not found", c.Formatter))
		}
		c.formatter = fmtr
	}
	return res
}

func (c *connWriter) SetFormatter(f LogFormatter) {
	c.formatter = f
}

// WriteMsg writes message in connection.
// If connection is down, try to re-connect.
func (c *connWriter) WriteMsg(lm *LogMsg) error {
	if lm.Level > c.Level {
		return nil
	}
	if c.needToConnectOnMsg() {
		err := c.connect()
		if err != nil {
			return err
		}
	}

	if c.ReconnectOnMsg {
		defer c.innerWriter.Close()
	}

	msg := c.formatter.Format(lm)

	_, err := c.lg.writeln(msg)
	if err != nil {
		return err
	}
	return nil
}

// Flush implementing method. empty.
func (c *connWriter) Flush() {
}

// Destroy destroy connection writer and close tcp listener.
func (c *connWriter) Destroy() {
	if c.innerWriter != nil {
		c.innerWriter.Close()
	}
}

func (c *connWriter) connect() error {
	if c.innerWriter != nil {
		c.innerWriter.Close()
		c.innerWriter = nil
	}

	conn, err := net.Dial(c.Net, c.Addr)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
	}

	c.innerWriter = conn
	c.lg = newLogWriter(conn)
	return nil
}

func (c *connWriter) needToConnectOnMsg() bool {
	if c.Reconnect {
		return true
	}

	if c.innerWriter == nil {
		return true
	}

	return c.ReconnectOnMsg
}

func init() {
	Register(AdapterConn, NewConn)
}
