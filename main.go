package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	setup "github.com/louisroyer/docker-setup/runtime"
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
