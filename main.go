package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"dotterel/engine"
	"dotterel/machine"
)

func main() {
	port := "/dev/ttyACM0" // Your Gemini PR device
	baud := 9600

	// Load your dictionary
	dict, err := engine.LoadDictionary("dict.json")
	if err != nil {
		log.Fatalf("Error loading dictionary: %v", err)
	}
	e := engine.NewEngine(dict)

	gemini := machine.NewGeminiPrMachine(port, baud, func(keys []string) {
		word := e.TranslateSteno(strings.Join(keys, " "))
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
