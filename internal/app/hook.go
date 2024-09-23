// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package app

// init & exit hooks iface
type Hook interface {
	String() []string
	RunInit() error
	RunExit() error
}

// init or exit hook iface
type HookSingle interface {
	String() []string
	Run() error
}

// init & exit hooks
type HookMulti struct {
	init HookSingle
	exit HookSingle
}

// Returns hooks info
func (hooks HookMulti) String() []string {
	r := []string{}
	for _, i := range hooks.init.String() {
		r = append(r, i)
	}
	for _, i := range hooks.exit.String() {
		r = append(r, i)
	}
	return r
}

// Runs init hook
func (hooks HookMulti) RunInit() error {
	return hooks.init.Run()
}

// Runs exit hook
func (hooks HookMulti) RunExit() error {
	return hooks.exit.Run()
}
