package global

import "time"

type SSClientChannels struct {
	Timeout      time.Duration
	ChanError    chan error
	ChanTCPReady chan bool
	ChanUDPReady chan bool
}
