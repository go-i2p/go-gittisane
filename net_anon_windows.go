// Copyright 2019 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

// This code is heavily inspired by the archived gofacebook/gracenet/net.go handler

//go:build windows

package graceful

import "net"

func GetListenerUnix(network string, addr net.Addr) (net.Listener, error) {
	return nil, nil
}
