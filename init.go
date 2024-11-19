package sshforward

import (
	"context"
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

type ForwardConfig_t struct {
	state forwardState_t
	wg    sync.WaitGroup

	sshConfig   *ssh.ClientConfig
	tunnelAddr  string
	eventNotify chan StateEvent_t

	client        *ssh.Client
	localListener *net.Listener

	ctx    context.Context
	cancel context.CancelFunc
}

func CreateForward() *ForwardConfig_t {
	ctx, cancel := context.WithCancel(context.Background())
	return &ForwardConfig_t{
		state: FORWARD_STATE_NONE,
		wg:    sync.WaitGroup{},

		sshConfig:   nil,
		tunnelAddr:  "",
		eventNotify: make(chan StateEvent_t, 10),

		client: nil,

		ctx:    ctx,
		cancel: cancel,
	}
}

func (f *ForwardConfig_t) ConfigTunnel(sshConfig *ssh.ClientConfig, tunnelAddr string, tunnelPort string) {
	f.sshConfig = sshConfig
	f.tunnelAddr = tunnelAddr + ":" + tunnelPort

	f.stateChange(FORWARD_STATE_CONFIGURED, "")
}
