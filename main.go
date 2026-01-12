// Copyright Louis Royer and the NextMN contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/nextmn/logrus-formatter/logger"

	"github.com/nextmn/docker-setup/internal/app"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

func main() {
	logger.Init("docker-setup")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	version := "Unknown version"
	if info, ok := debug.ReadBuildInfo(); ok {
		version = info.Main.Version
	}
	app := &cli.Command{
		Name:                  "docker-setup",
		Usage:                 "Docker-setup - Configure a container for networking",
		EnableShellCompletion: true,
		Authors: []any{
			"Louis Royer",
		},
		Version: version,
		Action: func(ctx context.Context, cmd *cli.Command) error {
			conf := app.NewConf()
			conf.RunInitHooks(ctx)
			if !conf.Oneshot() {
				defer func() {
					shCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 1*time.Second)
					defer cancel()
					conf.RunExitHooks(shCtx)
				}()
				<-ctx.Done()
			}
			return nil
		},
	}
	if err := app.Run(ctx, os.Args); err != nil {
		logrus.WithError(err).Fatal("Fatal error while running the application")
	}
}
