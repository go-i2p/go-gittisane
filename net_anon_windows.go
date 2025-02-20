// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// This code is heavily inspired by the archived gofacebook/gracenet/net.go handler

//go:build windows

package graceful

import "net"

func ResolveUnixAddr(network, address string) (net.Addr, error) {
	switch network {
	case "unix", "unixpacket":
		return net.ResolveUnixAddr(network, address)
	case "tcp", "tcp4", "tcp6":
		return net.ResolveTCPAddr(network, address)
	case "udp", "udp4", "udp6":
		return net.ResolveUDPAddr(network, address)
	default:
		return nil, net.UnknownNetworkError(network)
	}
}

func GetListenerUnixWrapper(network string, addr net.Addr) (net.Listener, error) {
	return net.Listen(network, addr.String())
}
