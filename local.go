package sshforward

import (
	"context"
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

func (f *ForwardConfig_t) remoteForwardToLocalListen(localPort, remotePort string) error {
	tunnelAddr := f.tunnelAddr

	client, err := ssh.Dial("tcp", tunnelAddr, f.sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial to remote SSH server. Error: %v", err)
	}

	f.client = client

	defer client.Close()

	localAddr := "localhost:" + localPort
	remoteAddr := "localhost:" + remotePort

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		return fmt.Errorf("Failed to listen on local port. Error: %v", err)
	}
	defer listener.Close()

	f.stateChange(FORWARD_STATE_SSH_CONNECTED, fmt.Sprintf("Listening on localhost:%s <-> %s:%s", localPort, tunnelAddr, remotePort))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f.wg.Done()

	for {
		localConn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return fmt.Errorf("Failed to accept connection")
			}
		}

		go f.forward(ctx, localConn, client, remoteAddr)
	}
}

func (f *ForwardConfig_t) forward(ctx context.Context, localConn net.Conn, sshConn *ssh.Client, remoteAddr string) {
	remoteConn, err := sshConn.DialContext(ctx, "tcp", remoteAddr)
	if err != nil {
		localConn.Close()
		return
	}

	go func() {
		io.Copy(remoteConn, localConn)
		remoteConn.Close()
	}()

	go func() {
		io.Copy(localConn, remoteConn)
		localConn.Close()
	}()
}
