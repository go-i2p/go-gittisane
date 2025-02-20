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
