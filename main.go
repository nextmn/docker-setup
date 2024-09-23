// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/nextmn/docker-setup/internal/app"
	"github.com/nextmn/docker-setup/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	logger.Init("docker-setup")
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()
	app := &cli.App{
		Name:                 "docker-setup",
		Usage:                "Configure a container for networking",
		EnableBashCompletion: true,
		Authors: []*cli.Author{
			{Name: "Louis Royer"},
		},
		Action: func(c *cli.Context) error {
			conf := app.NewConf()
			conf.RunInitHooks()
			if !conf.Oneshot() {
				select {
				case <-ctx.Done():
					break
				}
			}
			conf.RunExitHooks()
			return nil
		},
	}
	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}
