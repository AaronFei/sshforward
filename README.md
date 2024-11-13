[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# SSH Port Forwarding Package

A Go package that provides SSH port forwarding functionality with support for both local-to-remote and remote-to-local port forwarding.

## Features

- Local to remote port forwarding
- Remote to local port forwarding
- Event notification system
- Automatic retry on connection failure
- Concurrent connection handling
- State management

## Installation

```bash
go get github.com/yourusername/sshforward
```

## Usage

### Basic Setup

```go
import (
    "github.com/yourusername/sshforward"
    "golang.org/x/crypto/ssh"
)

// Create a new forward instance
forward := sshforward.CreateForward()

// Configure SSH client
sshConfig := &ssh.ClientConfig{
    User: "username",
    Auth: []ssh.AuthMethod{
        ssh.Password("password"),
        // or use ssh.PublicKeys(...)
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
}

// Configure the tunnel
forward.ConfigTunnel(sshConfig, "example.com", "22")
```

### Starting Port Forwarding

#### Local to Remote Port Forwarding
```go
// Forward local port 8080 to remote port 80
forward.Service(sshforward.FORWARD_TYPE_LOCAL_TO_REMOTE_LISTEN, "80", "8080")

// Wait for the forwarding to be ready
forward.Wait()

// Now the port forwarding is ready to use
```

#### Remote to Local Port Forwarding
```go
// Forward remote port 80 to local port 8080
forward.Service(sshforward.FORWARD_TYPE_REMOTE_TO_LOCAL_LISTEN, "80", "8080")

// Wait for the forwarding to be ready
forward.Wait()

// Now the port forwarding is ready to use
```

### Event Monitoring

```go
// Get the event notification channel
eventChan := forward.EventNotifyChannel()

// Monitor events
go func() {
    for event := range eventChan {
        fmt.Printf("State: %s, Time: %s, Message: %s\n",
            event.State, event.T.Format(time.RFC3339), event.Msg)
    }
}()
```

## States

The forwarding service can be in following states:

- `NONE`: Initial state
- `CONFIGURED`: SSH configuration completed
- `STARTING`: Service is starting
- `SSH_CONNECTED`: SSH connection established and ready for use
- `STOPPED`: Service stopped
- `SKIP`: Service skipped (e.g., port unavailable)
- `ERROR`: Error occurred
- `RETRY`: Service is retrying after failure

## Complete Example

```go
package main

import (
    "fmt"
    "time"
    "github.com/yourusername/sshforward"
    "golang.org/x/crypto/ssh"
)

func main() {
    // Create and configure forward
    forward := sshforward.CreateForward()
    
    sshConfig := &ssh.ClientConfig{
        User: "username",
        Auth: []ssh.AuthMethod{
            ssh.Password("password"),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    
    forward.ConfigTunnel(sshConfig, "example.com", "22")
    
    // Monitor events
    go func() {
        for event := range forward.EventNotifyChannel() {
            fmt.Printf("[%s] %s: %s\n", 
                event.T.Format(time.RFC3339),
                event.State,
                event.Msg)
        }
    }()
    
    // Start forwarding
    forward.Service(sshforward.FORWARD_TYPE_LOCAL_TO_REMOTE_LISTEN, "80", "8080")
    
    // Wait for forwarding to be ready
    forward.Wait()
    
    fmt.Println("Port forwarding is now ready to use!")
    
    // Keep the program running
    select {}
}
```

## Error Handling

The service automatically retries on connection failures with a 10-second delay. Error messages are sent through the event notification channel.

## Notes

- Ensure the target ports are available before starting the service
- The service runs in a separate goroutine
- The event channel has a buffer size of 10
- For remote to local forwarding, the local port availability is checked before starting
- `Wait()` blocks until the SSH forwarding is ready to use
- The service will continue running in the background after `Wait()` returns

