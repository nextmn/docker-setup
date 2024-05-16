// Copyright 2023 Louis Royer and docker-setup contributors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.
// SPDX-License-Identifier: MIT
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	setup "github.com/nextmn/docker-setup/runtime"
)

// Initialize signals handling
func initSignals(conf setup.Conf) {
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)
	func(_ os.Signal) {}(<-cancelChan)
	conf.RunExitHooks()
	os.Exit(0)
}

// Print the configuration, then startup
func main() {
	log.SetPrefix("[docker-setup]")
	conf := setup.NewConf()
	conf.Log()
	go initSignals(conf)
	conf.RunInitHooks()
	if !conf.Oneshot() {
		select {}
	}
}
