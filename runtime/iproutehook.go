// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package setup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Runs iptables
func runIPRoute(args ...string) error {
	r := []string{"route"}
	r = append(r, args...)
	cmd := exec.Command("ip", r...)
	return cmd.Run()
}

// Create init & exit hooks
func NewIPRouteHooks(nameInit string, envInit string, nameExit string, envExit string) HookMulti {
	return HookMulti{
		init: NewIPRouteHook(nameInit, envInit),
		exit: NewIPRouteHook(nameExit, envExit),
	}
}

// IP Route Hook
type IPRouteHook struct {
	name   string // name of the hook
	env    string // environment variable
	isset  bool   // false when no cmd is set
	routes [][]string
}

// Create an IP Route Hook
func NewIPRouteHook(name string, env string) IPRouteHook {
	hook := IPRouteHook{
		name: name,
		env:  env,
	}
	if ifaces, ok := os.LookupEnv(hook.env); !ok {
		hook.isset = false
	} else {
		rl := strings.Split(ifaces, "\n")

		hook.routes = make([][]string, len(rl))
		for i, route := range rl {
			hook.routes[i] = strings.Split(strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(route), "-")), " ")
		}
		hook.isset = true
	}
	return hook
}

// Run the hook if it is set
func (hook IPRouteHook) Run() error {
	if !hook.isset {
		return nil
	}
	for _, r := range hook.routes {
		r = append(r, "proto", "static")
		if err := runIPRoute(r...); err != nil {
			return err
		}
	}
	return nil
}

// Returns hook information in an human format
func (hook IPRouteHook) String() []string {
	r := []string{}
	if !hook.isset {
		return []string{fmt.Sprintf("%s hook ($%s) is not set.", hook.name, hook.env)}
	}
	for i, route := range hook.routes {
		args := ""
		for _, cmd := range route {
			args += fmt.Sprintf(", %s", cmd)
		}
		r = append(r, fmt.Sprintf("%s #%d hook is set with : [ip, route%s, proto, static]", hook.name, i, args))
	}
	return r
}
