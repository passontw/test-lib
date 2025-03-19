package types

import "time"

type HeartbeatStatus struct {
	ServerId   string
	UpdateTime time.Time
}
