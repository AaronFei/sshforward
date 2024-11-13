package sshforward

import "time"

type forwardState_t int

const (
	FORWARD_STATE_NONE forwardState_t = iota
	FORWARD_STATE_CONFIGURED
	FORWARD_STATE_STARTING
	FORWARD_STATE_SSH_CONNECTED
	FORWARD_STATE_STOPPED
	FORWARD_STATE_SKIP
	FORWARD_STATE_ERROR
	FORWARD_STATE_RETRY
)

var forwardStateString = map[forwardState_t]string{
	FORWARD_STATE_NONE:          "NONE",
	FORWARD_STATE_CONFIGURED:    "CONFIGURED",
	FORWARD_STATE_STARTING:      "STARTING",
	FORWARD_STATE_SSH_CONNECTED: "SSH_CONNECTED",
	FORWARD_STATE_STOPPED:       "STOPPED",
	FORWARD_STATE_SKIP:          "SKIP",
	FORWARD_STATE_ERROR:         "ERROR",
	FORWARD_STATE_RETRY:         "RETRY",
}

type StateEvent_t struct {
	State string
	T     time.Time
	Msg   string
}

func (f *forwardConfig_t) stateChange(newState forwardState_t, msg string) {
	f.state = newState
	select {
	case f.eventNotify <- StateEvent_t{State: forwardStateString[newState], T: time.Now(), Msg: msg}:
	default:
	}
}

func (f *forwardConfig_t) EventNotifyChannel() <-chan StateEvent_t {
	return f.eventNotify
}
