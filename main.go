package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"skyfall.mcafee.int/intel/logger"
)

var (
	log = logger.NewLogger("backendmock")
)

func main() {
	fmt.Println("Mock GTI server")
	fmt.Println("Git hash: ", githash)
	fmt.Println("Mock Backend Server")
	fmt.Println("Testsetupserver Server")
	fmt.Println("Build time: ", buildstamp)
	fmt.Println("Build number: ", buildnumber)
	getConfig()

	if ok := Serve(); ok {
		sig := make(chan os.Signal)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	forever:
		for {
			select {
			case s := <-sig:
				fmt.Println("Signal ", s, " received, stopping")
				break forever
			}
		}
		fmt.Println("Server stopped")
	}

}
