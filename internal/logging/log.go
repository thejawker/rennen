package logging

import (
	"fmt"
	"log"
	"os"
)

func SetupLogging() func() error {
	logFile, err := os.Create("ren_debug.log")
	if err != nil {
		fmt.Printf("Error creating log file: %v\n", err)
		return func() error { return nil }
	}
	log.SetOutput(logFile)

	// defer logFile.Close()
	return logFile.Close
}
