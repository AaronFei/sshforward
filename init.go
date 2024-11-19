package sshforward

import (
	"sync"

	"golang.org/x/crypto/ssh"
)

type ForwardConfig_t struct {
	state forwardState_t
	wg    sync.WaitGroup

	sshConfig   *ssh.ClientConfig
	tunnelAddr  string
	eventNotify chan StateEvent_t

	client *ssh.Client
}

func CreateForward() *ForwardConfig_t {
	return &ForwardConfig_t{
		state: FORWARD_STATE_NONE,
		wg:    sync.WaitGroup{},

		sshConfig:   nil,
		tunnelAddr:  "",
		eventNotify: make(chan StateEvent_t, 10),

		client: nil,
	}
}

func (f *ForwardConfig_t) ConfigTunnel(sshConfig *ssh.ClientConfig, tunnelAddr string, tunnelPort string) {
	f.sshConfig = sshConfig
	f.tunnelAddr = tunnelAddr + ":" + tunnelPort

	f.stateChange(FORWARD_STATE_CONFIGURED, "")
}
