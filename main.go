// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"sten/config"
	"sten/engine"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	e := engine.NewEngine(cfg)

	go e.Run()

	fmt.Println("[sten] Running. Press Ctrl+C to quit.")

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("\n[sten] Quit with Ctrl+C")
}
