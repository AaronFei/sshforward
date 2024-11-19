package sshforward

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"

	"golang.org/x/crypto/ssh"
)

func (f *ForwardConfig_t) localForwardToRemoteListen(remotePort, localPort string) error {
	tunnelAddr := f.tunnelAddr

	client, err := ssh.Dial("tcp", tunnelAddr, f.sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to dial to remote SSH server. Error: %v", err)
	}

	f.client = client

	defer client.Close()

	tunnelAddrWithoutPort := tunnelAddr[:strings.LastIndex(tunnelAddr, ":")]
	listener, err := client.Listen("tcp", fmt.Sprintf("localhost:%s", remotePort))
	if err != nil {
		return fmt.Errorf("Failed to listen on remote port. Error: %v", err)
	}

	f.localListener = &listener

	defer listener.Close()

	f.stateChange(FORWARD_STATE_SSH_CONNECTED, fmt.Sprintf("Listening on %s:%s <-> localhost:%s", tunnelAddrWithoutPort, remotePort, localPort))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f.wg.Done()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if err == io.EOF {
				return nil
			} else {
				return fmt.Errorf("Failed to accept connection")
			}
		}

		go f.handleConnection(ctx, conn, localPort)
	}

}

func (f *ForwardConfig_t) handleConnection(ctx context.Context, conn net.Conn, localPort string) {
	defer conn.Close()

	dialer := net.Dialer{}

	localConn, err := dialer.DialContext(ctx, "tcp", fmt.Sprintf("localhost:%s", localPort))
	if err != nil {
		return
	}
	defer localConn.Close()

	go io.Copy(localConn, conn)

	io.Copy(conn, localConn)
}
