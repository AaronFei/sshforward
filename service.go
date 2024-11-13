package sshforward

import (
	"fmt"
	"net"
	"time"
)

type ForwardType_t int

const (
	FORWARD_TYPE_LOCAL_TO_REMOTE_LISTEN ForwardType_t = iota
	FORWARD_TYPE_REMOTE_TO_LOCAL_LISTEN
)

func (f *ForwardConfig_t) Wait() {
	f.wg.Wait()
}

func (f *ForwardConfig_t) Service(t ForwardType_t, remotePort, localPort string) {
	if f.state != FORWARD_STATE_CONFIGURED || f.state == FORWARD_STATE_STOPPED {
		f.stateChange(FORWARD_STATE_ERROR, fmt.Sprintf("Not at configured state or stopped state"))
		return
	}

	f.wg.Add(1)

	go f.forwardService(t, remotePort, localPort)
}

func (f *ForwardConfig_t) forwardService(t ForwardType_t, remotePort, localPort string) {
	f.stateChange(FORWARD_STATE_STARTING, fmt.Sprintf("Forward service started. Type: %d, Remote port: %s, Local port: %s", t, remotePort, localPort))

	for {
		var err error = nil

		if t == FORWARD_TYPE_REMOTE_TO_LOCAL_LISTEN {
			if !localPortAvailable(localPort) {
				f.stateChange(FORWARD_STATE_SKIP, fmt.Sprintf("Local port %s is not available for listening", localPort))
				return
			}
		}

		switch t {
		case FORWARD_TYPE_LOCAL_TO_REMOTE_LISTEN:
			err = f.localForwardToRemoteListen(remotePort, localPort)
		case FORWARD_TYPE_REMOTE_TO_LOCAL_LISTEN:
			err = f.remoteForwardToLocalListen(localPort, remotePort)
		}

		if err == nil {
			f.stateChange(FORWARD_STATE_STOPPED, "")
			break
		}

		f.stateChange(FORWARD_STATE_RETRY, fmt.Sprintf("Forward service failed. Retry in 10 seconds. Error: %v", err))
		time.Sleep(10 * time.Second)
	}
}

func localPortAvailable(localPort string) bool {
	l, err := net.Listen("tcp", "localhost:"+localPort)

	if l != nil {
		defer l.Close()
	}

	if err != nil {
		return false
	}

	return true
}
