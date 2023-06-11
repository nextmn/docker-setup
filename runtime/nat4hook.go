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
func runIP4Tables(args ...string) error {
	cmd := exec.Command("iptables", args...)
	return cmd.Run()
}

// IPv4 NAT Hook
type Nat4Hook struct {
	name   string // name of the hook
	env    string // environment variable
	isset  bool   // false when no cmd is set
	ifaces []string
}

// Creates a new IPv4 NAT Hook
func NewNat4Hooks() Nat4Hook {
	hook := Nat4Hook{
		name: "IPv4 NAT",
		env:  "NAT4_IFACES",
	}
	if ifaces, ok := os.LookupEnv(hook.env); !ok {
		hook.isset = false
	} else {
		ifaces_l := strings.Split(ifaces, "\n")
		hook.ifaces = make([]string, len(ifaces_l))
		for i, iface := range ifaces_l {
			hook.ifaces[i] = strings.TrimSpace(strings.TrimLeft(strings.TrimSpace(iface), "-"))
		}
		hook.isset = true
	}

	return hook
}

// Runs IPv4 NAT init hook
func (hook Nat4Hook) RunInit() error {
	if !hook.isset {
		return nil
	}
	if err := runIP4Tables("-I", "FORWARD", "-j", "ACCEPT"); err != nil {
		return err
	}
	for _, iface := range hook.ifaces {
		if err := runIP4Tables("-t", "nat", "-A", "POSTROUTING", "-o", iface, "-j", "MASQUERADE"); err != nil {
			return err
		}
	}
	return nil
}

// Runs IPv4 NAT exit hook
func (hook Nat4Hook) RunExit() error {
	if !hook.isset {
		return nil
	}
	errcount := 0
	var lasterr error
	for _, iface := range hook.ifaces {
		// if there is an error, we continue: will return at the end
		if err := runIP4Tables("-t", "nat", "-D", "POSTROUTING", "-o", iface, "-j", "MASQUERADE"); err != nil {
			errcount++
			lasterr = err
		}
	}
	if errcount == 1 {
		return lasterr
	} else if errcount > 1 {
		return fmt.Errorf("%i iptable errors", errcount)
	}
	return nil
}

// Returns IPv4 NAT hook infos
func (hook Nat4Hook) String() []string {
	r := []string{}
	if !hook.isset {
		r = append(r, fmt.Sprintf("%s hook is not set ($%s (list)).", hook.name, hook.env))
		return r
	}
	for _, i := range hook.ifaces {
		r = append(r, fmt.Sprintf("%s is set for interface %s", hook.name, i))
	}
	return r
}
