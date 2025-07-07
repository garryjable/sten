// Copyright (c) 2025 Garrett Jennings.
// See LICENSE.txt for details.
// This file is part of GPlover.
// GPlover is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.
// This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.
// You should have received a copy of the GNU General Public License along with this program. If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gplover/config"
	"gplover/dictionary"
	"gplover/engine"
	"gplover/machine"
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
	e := engine.NewEngine(dict)

	gemini := machine.NewGeminiPrMachine(cfg.Port, cfg.Baud, func(keys []string) {
		word := e.TranslateSteno(keys)
		fmt.Print(word + " ")
	})

	// Start machine capture
	err = gemini.StartCapture()
	if err != nil {
		log.Fatalf("Failed to start Gemini PR machine: %v", err)
	}
	defer gemini.StopCapture()

	fmt.Println("Dotterel now running. Press Ctrl+C to quit.")

	// Handle Ctrl+C
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	<-sigs
	fmt.Println("\n[dotterel] Quit with Ctrl+C")
}
