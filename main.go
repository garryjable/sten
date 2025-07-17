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
	"sten/dictionary"
	"sten/engine"
	"sten/machine"
	"sten/output"
	"sten/stroke"
	"sten/translator"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Load your dictionary
	dict, longestOutline, err := dictionary.LoadDictionaries("dictionaries")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}

	// Create virtual Output
	out, err := output.NewVirtualOutput()
	if err != nil {
		log.Fatalf("Failed to init virtual keyboard: %v", err)
	}

	e := engine.NewEngine(out)

	t := translator.NewTranslator(dict, longestOutline)

	gemini := machine.NewGeminiPrMachine(cfg.Port, cfg.Baud, func(stroke *stroke.Stroke) {
		// word := t.translate(stroke)
		translation := t.Translate(stroke.Steno())
		e.Execute(translation)

	})
	// Start machine capture
	err = gemini.StartCapture()
	if err != nil {
		log.Fatalf("Failed to start Gemini PR machine: %v", err)
	}
	defer gemini.StopCapture()

	fmt.Println("[sten] Running. Press Ctrl+C to quit.")

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("\n[sten] Quit with Ctrl+C")
}
