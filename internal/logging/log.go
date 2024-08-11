package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

func SetupLogging(level *string) func() error {
	if level == nil || *level == "none" {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

		// preventing logging anything to the console
		log.SetOutput(io.Discard)

		return func() error { return nil }
	}

	logFile, err := os.Create("ren.log")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return func() error { return nil }
	}

	log.SetOutput(logFile)

	// defer logFile.Close()
	return logFile.Close
}
