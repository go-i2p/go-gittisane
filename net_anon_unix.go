// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// This code is heavily inspired by the archived gofacebook/gracenet/net.go handler

//go:build !windows

package graceful

import (
	"fmt"
	"net"
)

func ResolveUnixAddr(network, address string) (net.Addr, error) {
	switch network {
	case "unix", "unixpacket":
		return net.ResolveUnixAddr(network, address)
	default:
		return nil, fmt.Errorf("unknown network type %s", network)
	}
}

func GetListenerUnixWrapper(network string, addr net.Addr) (net.Listener, error) {
	switch addr.(type) {
	case *net.UnixAddr:
		return GetListenerUnix(network, addr.(*net.UnixAddr))
	default:
		return nil, fmt.Errorf("unknown address type %T", addr)
	}

}
