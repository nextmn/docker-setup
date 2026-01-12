// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"context"
)

// init & exit hooks iface
type Hook interface {
	String() []string
	RunInit(context.Context) error
	RunExit(context.Context) error
}

// init or exit hook iface
type HookSingle interface {
	String() []string
	Run(context.Context) error
}

// init & exit hooks
type HookMulti struct {
	init HookSingle
	exit HookSingle
}

// Returns hooks info
func (hooks HookMulti) String() []string {
	r := []string{}
	r = append(r, hooks.init.String()...)
	r = append(r, hooks.exit.String()...)
	return r
}

// Runs init hook
func (hooks HookMulti) RunInit(ctx context.Context) error {
	return hooks.init.Run(ctx)
}

// Runs exit hook
func (hooks HookMulti) RunExit(ctx context.Context) error {
	return hooks.exit.Run(ctx)
}
