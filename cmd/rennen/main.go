package main

import (
	"flag"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thejawker/rennen/internal/config"
	"github.com/thejawker/rennen/internal/logging"
	"github.com/thejawker/rennen/internal/model"
	"github.com/thejawker/rennen/internal/process"
	"log"
	"os"
)

func main() {
	configPath := flag.String("config", "./ren.json", "path to config file")
	showVersion := flag.Bool("version", false, "show version information")
	verbosityLevel := flag.String("logging", "none", "logs to ./ren.log verbosity level: none, all")
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("Rennen version %s\n", getVersion())
		return
	}

	// Setup logging
	closeLogger := logging.SetupLogging(verbosityLevel)

	// Defer and handle the error
	defer func() {
		if err := closeLogger(); err != nil {
			fmt.Printf("error closing log file: %w", err)
		}
	}()

	// Load configuration
	cfg, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize processes
	processes, err := process.InitializeFromConfig(cfg.Processes)
	if err != nil {
		log.Fatalf("Error initializing processes: %v", err)
	}

	// Create and initialize the model
	m := model.New(processes)

	// Create and start the Bubble Tea program
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	fmt.Println("\nwoah that was cool i guess, all is stopped now though")
}

func getVersion() string {
	// get from the ./VERSION file
	file, err := os.ReadFile("./VERSION")
	if err != nil {
		return "unknown"
	}

	return string(file)
}

func loadConfig(configPath string) (*config.Config, error) {
	if configPath == "" {
		configPath = "ren.json"
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return nil, fmt.Errorf("error loading config from %s: %w", configPath, err)
	}

	return cfg, nil
}
