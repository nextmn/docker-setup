// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/sirupsen/logrus"
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
func (hook UserHook) Run(ctx context.Context) error {
	if !hook.isset {
		return nil
	}
	cmd := exec.CommandContext(ctx, hook.cmd, hook.args...)
	cmd.Env = []string{}
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		logrus.WithError(err).Error("Error running user hook")
		return err
	}
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		logrus.WithFields(logrus.Fields{
			"hook-name": hook.name,
			"message":   m,
		}).Info("Hook output")
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
