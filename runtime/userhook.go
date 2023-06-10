// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package setup

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

// Returns a list of arguments for a command defined in from environment variables.
func get_args(env_var_prefix string) []string {
	hook_args := []string{}
	for num_arg := 0; ; num_arg++ {
		if arg, isarg := os.LookupEnv(fmt.Sprintf("%s%d", env_var_prefix, num_arg)); !isarg {
			return hook_args
		} else {
			hook_args = append(hook_args, arg)
		}
	}
}

// Create init & exit hooks
func NewUserHooks(nameInit string, envInit string, nameExit string, envExit string) HookMulti {
	return HookMulti{
		init: NewUserHook(nameInit, envInit),
		exit: NewUserHook(nameExit, envExit),
	}
}

// Hook
type UserHook struct {
	name  string   // name of the hook
	isset bool     // false when no cmd is set
	env   string   // environment variable used to retrieve configuration
	cmd   string   // command
	args  []string // argument list
}

// Creates a new Hook from environment variables
func NewUserHook(name string, env string) UserHook {
	hook := UserHook{
		env:  env,
		name: name,
	}
	if cmd, ok := os.LookupEnv(env); !ok {
		hook.isset = false
	} else {
		hook.cmd = cmd
		hook.isset = true
	}
	hook.args = get_args(fmt.Sprintf("%s_", env))
	return hook
}

// Run the hook if it is set
func (hook UserHook) Run() error {
	if !hook.isset {
		return nil
	}
	cmd := exec.Command(hook.cmd, hook.args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		log.Printf("[%s hook] %s", hook.name, m)
	}
	return cmd.Wait()
}

// Returns hook information in an human format
func (hook UserHook) String() []string {
	if !hook.isset {
		return []string{fmt.Sprintf("%s hook ($%s) is not set.", hook.name, hook.env)}
	}
	args := ""
	for _, a := range hook.args {
		args += fmt.Sprintf(", %s", a)
	}
	return []string{fmt.Sprintf("%s hook is set with : [%s%s]", hook.name, hook.cmd, args)}
}
