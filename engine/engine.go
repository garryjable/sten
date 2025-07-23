// Copyright (c) 2025 Garrett Jennings.
// This File is part of sten. Sten is free software under GPLv3 .
// See LICENSE.txt for details.

package engine

import (
	"log"
	"sten/config"
	"sten/dictionary"
	"sten/machine"
	"sten/output"
	"sten/translator"
)

type Engine struct {
	cfg        *config.Config
	machine    machine.Machine
	translator *translator.Translator
	output     output.OutputService
}

func NewEngine(cfg *config.Config) *Engine {
	// Load your dictionary
	dict, longestOutline, err := dictionary.LoadDictionaries("dictionaries")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}

	var o output.OutputService
	var m machine.Machine
	if cfg.Machine == "geminipr" {
		m = machine.NewGeminiPrMachine(cfg.Port, cfg.Baud)
	} else {
		log.Fatalf("Unknown machine type: %v", cfg.Machine)
	}
	t := translator.NewTranslator(dict, longestOutline, m.Strokes())
	if cfg.Dev == true {
		o = output.NewDevOutputService(t.Out())
	} else {
		log.Fatalf("Non-dev mode output not implemented!")
	}

	e := &Engine{
		cfg:        cfg,
		machine:    m,
		output:     o,
		translator: t,
	}
	return e
}

func (e *Engine) Run() {
	// Start machine capture
	go e.machine.StartCapture()
	go e.translator.Run()
	go e.output.Run()
	defer e.machine.StopCapture()
}
