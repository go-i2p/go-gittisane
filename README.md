# Go-Gittisane

A soft-fork of Gitea with built-in support for running as an I2P (Invisible Internet Project) service. This project provides network implementation files and CI configuration to build Gitea with I2P support.

## What is Go-Gittisane?

Go-Gittisane enables you to run the Gitea git server with anonymity features through the I2P network. All HTTP traffic - both incoming requests to your Gitea instance and outgoing requests from it - are routed through I2P, providing privacy and censorship resistance.

## How it Works

### Technical Explanation

Go-Gittisane leverages Gitea's modular network architecture to replace the standard TCP/IP networking with I2P connectivity:

1. **Network Module Substitution**: Gitea encapsulates its network operations in the `graceful` package. By providing alternative implementations of key functions like `GetListener()`, we can redirect all network traffic through I2P.

2. **I2P Integration**: The core modification is in `net_anon.go`, which uses the `github.com/go-i2p/onramp` library to establish I2P connectivity. This file:
   - Creates an I2P "Garlic" router connection
   - Implements a custom `GetListener` function that returns I2P listeners 
   - Configures Go's HTTP client to use I2P for outbound connections

3. **Platform-Specific Handling**: Unix socket connections (used for local IPC) are handled normally since they don't present anonymity risks. These implementations are in `net_anon_unix.go` and `net_anon_windows.go`.

### Automated Builds

This repository contains a GitHub Actions workflow that:
1. Checks for new Gitea releases daily
2. Downloads the Gitea source code for the latest release
3. Applies our I2P network modifications
4. Builds binaries for Linux, Windows, and macOS
5. Creates a new release with these modified binaries

## Deployment Options

### Using Pre-built Binaries

1. Download the latest release for your platform from the [Releases page](https://github.com/go-i2p/go-gittisane/releases)
2. Ensure you have I2P router running with SAM enabled on port 7656
3. Run the gittisane binary(it has identical options as gitea)

## Implementation Details

The core networking modification looks like this:

```go
// This file gets copied to modules/graceful/net_anon.go during build
package graceful

import (
    "net"
    "net/http"

    "github.com/go-i2p/onramp"
)

// Set up I2P connectivity
var garlic, i2perr = onramp.NewGarlic("gitea-anon", "127.0.0.1:7656", onramp.OPT_DEFAULTS)

// Custom implementation of GetListener for I2P
func I2PGetListener(network, address string) (net.Listener, error) {
    defer GetManager().InformCleanup()
    switch network {
    case "tcp", "tcp4", "tcp6", "i2p", "i2pt":
        return garlic.Listen()
    case "unix", "unixpacket":
        // Unix sockets handled normally
        unixAddr, err := ResolveUnixAddr(network, address)
        if err != nil {
            return nil, err
        }
        return GetListenerUnixWrapper(network, unixAddr)
    default:
        return nil, net.UnknownNetworkError(network)
    }
}

// Initialize everything at runtime
func init() {
    if i2perr != nil {
        panic(i2perr)
    }
    GetListener = I2PGetListener
    httpClient := &http.Client{
        Transport: &http.Transport{
            Dial: garlic.Dial,
        },
    }

    http.DefaultClient = httpClient
    http.DefaultTransport = httpClient.Transport
}
```

## Caveats and Limitations

While HTTP traffic is anonymized, other types of connections might not be:

1. **SMTP**: Email sending from Gitea is not automatically routed through I2P
2. **SSH**: Git operations using SSH are not automatically anonymized
3. **External Services**: Webhooks and other external connections will use I2P, but services must support I2P addresses

These limitations can be addressed with additional configuration but are beyond the scope of the default implementation.

## Requirements

- Go 1.21 or later (for building from source)
- Running I2P router with SAM API enabled on port 7656
- Standard Gitea requirements (database, etc.)

## License

Both this modification and Gitea are licensed under the MIT license.
- See [LICENSE](LICENSE) for the license covering the files in this repository
- See [LICENSE-gitea.md](LICENSE-gitea.md) for the Gitea license

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.