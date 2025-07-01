package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"dotterel/engine"

	"golang.org/x/term"
)

func main() {
	// port := "/dev/ttyUSB0" // or COM3 on Windows
	// baud := 9600

	// gemini := machine.NewGeminiPrMachine(port, baud, func(keys []string) {
	// 	fmt.Printf("Stroke keys: %v\n", keys)
	// })

	// err := gemini.StartCapture()
	// if err != nil {
	// 	log.Fatalf("Failed to start Gemini PR machine: %v", err)
	// }
	// defer gemini.StopCapture()

	dict, err := engine.LoadDictionary("dict.json")
	if err != nil {
		fmt.Println("Error loading dictionary:", err)
		return
	}

	e := engine.NewEngine(dict)

	fmt.Println("Dotterel running in raw mode. Press Ctrl+C to quit.")

	fd := int(os.Stdin.Fd())
	oldState, err := term.MakeRaw(fd)
	if err != nil {
		fmt.Println("Failed to enter raw mode:", err)
		return
	}
	defer term.Restore(fd, oldState)

	// Set up Ctrl+C signal handling
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		term.Restore(fd, oldState)
		fmt.Println("\n[dotterel] Quit with Ctrl+C")
		os.Exit(0)
	}()

	var buffer []rune
	stdin := os.NewFile(uintptr(fd), "/dev/stdin")

	for {
		b := make([]byte, 1)
		_, err := stdin.Read(b)
		if err != nil {
			fmt.Println("Read error:", err)
			break
		}

		switch b[0] {
		case 0x03: // Ctrl+C ASCII code
			term.Restore(fd, oldState)
			fmt.Println("\n[dotterel] Quit with Ctrl+C")
			os.Exit(0)
		case 0x00:
			continue
		case ' ':
			if len(buffer) > 0 {
				stroke := string(buffer)
				word := e.TranslateSteno(stroke)
				fmt.Print(word + " ")
				buffer = nil
			}
		default:
			buffer = append(buffer, rune(b[0]))
		}
	}
}
