package sshforward

import (
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
)

type ForwardConfig_t struct {
	state forwardState_t
	wg    sync.WaitGroup

	sshConfig   *ssh.ClientConfig
	TunnelAddr  string
	EventNotify chan StateEvent_t
}

func CreateForward() *ForwardConfig_t {
	return &ForwardConfig_t{
		state: FORWARD_STATE_NONE,
		wg:    sync.WaitGroup{},

		sshConfig:   nil,
		TunnelAddr:  "",
		EventNotify: make(chan StateEvent_t, 10),
	}
}

func (f *ForwardConfig_t) ConfigTunnel(sshConfig *ssh.ClientConfig, tunnelAddr string, tunnelPort string) {
	f.sshConfig = sshConfig
	f.TunnelAddr = tunnelAddr + ":" + tunnelPort

	f.stateChange(FORWARD_STATE_CONFIGURED, "")
}
