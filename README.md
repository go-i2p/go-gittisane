# go-gittisane
A soft-fork of gitea with support for running as an I2P service. Just the mod and the CI files.

How it works:
=============

Very simply, this uses Github CI to continuously build an I2P-only version of gitea based on the latest release of gitea at all times.
We can do this without requiring a patch to the gitea source code.
This is because gitea encapsulates it's "Listening" and "Dialing" into functions, which can easily be substituted for alternative versions.
For instance, the network listener is set up by a function, `graceful.GetListener() (net.Listener, error)` in the file `modules/graceful/server.go`. 
The default implementation of the `GetListener() (net.Listener, error)` function, `DefaultGetListener() (net.Listener, error)` is defined in the `modules/graceful/net_unix.go` for Unix-like systems and `modules/graceful/net_windows.go` for Windows-like systems.
A developer who wishes to "Mod" gitea to listen on another kind of connection do so by creating a new file which implements a `GetListener() (net.Listener, error)` function using an alternate listener implementation.

On the client side, the same thing is possible because Go allows you to substitute the underlying transports used for the default HTTP Client.
So, in the absence of overriding settings, we can configure it to use SAMv3 to build HTTP connections by default using the same keys as the hidden service.
Effectively this is like a "Bidirectional HTTP" tunnel in Java's Hidden Service Manager.

Finally, if you need to include additional libraries, run `go mod tidy` in the root of the gitea checkout to include them.

Here is a complete working example mod:

```Go
// copy this file to modules/graceful/net_anon.go before building gitea
package graceful

import (
	"net"
	"net/http"

	"github.com/go-i2p/onramp"
)

// First, make sure that the onramp.Garlic API is set up:
var garlic, i2perr = onramp.NewGarlic("gitea-anon", "127.0.0.1:7656", onramp.OPT_DEFAULTS)

// This implements the GetListener function for I2P. Note the exemption for Unix sockets.
func I2PGetListener(network, address string) (net.Listener, error) {
	// Add a deferral to say that we've tried to grab a listener
	defer GetManager().InformCleanup()
	switch network {
	case "tcp", "tcp4", "tcp6", "i2p", "i2pt":
		return garlic.Listen()
	case "unix", "unixpacket":
		// I2P isn't really a replacement for the stuff you use Unix sockets for and it's also not an anonymity risk, so treat them normally
		unixAddr, err := net.ResolveUnixAddr(network, address)
		if err != nil {
			return nil, err
		}
		return GetListenerUnix(network, unixAddr)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

// We use `init() to ensure that the I2P Listeners and Dialers are correctly placed at runtime`
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