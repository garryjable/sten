// Copyright (c) 2025 Garrett Jennings.
// This File is part of gplover. Gplover is free software under GPLv3 .
// See LICENSE.txt for details.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gplover/config"
	"gplover/dictionary"
	"gplover/machine"
	"gplover/output"
	"gplover/stroke"
	"gplover/translator"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load your dictionary
	dict, err := dictionary.LoadDictionaries("dictionaries")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}
	//e := engine.NewEngine(dict)

	// Create virtual Output
	out, err := output.NewVirtualOutput()
	if err != nil {
		log.Fatalf("Failed to init virtual keyboard: %v", err)
	}
	defer out.Close()

	t := translator.NewTranslator(dict, 1000)

	gemini := machine.NewGeminiPrMachine(cfg.Port, cfg.Baud, func(stroke *stroke.Stroke) {
		// word := t.translate(stroke)
		translation := t.Translate(stroke)
		_ = out.TypeString(translation.English + " ")

	})
	// Start machine capture
	err = gemini.StartCapture()
	if err != nil {
		log.Fatalf("Failed to start Gemini PR machine: %v", err)
	}
	defer gemini.StopCapture()

	fmt.Println("Gplover now running. Press Ctrl+C to quit.")

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("\n[gplover] Quit with Ctrl+C")
}
