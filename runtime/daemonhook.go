// Copyright 2024 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package setup

type DaemonHook struct {
	exit *chan bool
}

func NewDaemonHook(exit *chan bool) DaemonHook {
	return DaemonHook{
		exit: exit,
	}
}

func (hook DaemonHook) String() []string {
	return []string{}
}

func (hook DaemonHook) RunInit() error {
	// start the entrypoint and daemon goroutines
	return nil
}

func (hook DaemonHook) RunExit() error {
	// stop  entrypoint and daemon goroutines
	return nil
}
