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
)

var (
	version = "0.0.3" // This should be updated with each release
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "", "path to config file (default: ./ren.json)")
	showVersion := flag.Bool("version", false, "show version information")
	verbose := flag.Bool("verbose", false, "enable verbose logging")
	flag.Parse()

	// Show version and exit if requested
	if *showVersion {
		fmt.Printf("Rennen version %s\n", version)
		return
	}

	// Setup logging
	closeLogger := logging.SetupLogging()

	// defer and handle the error
	defer func() {
		if err := closeLogger(); err != nil {
			fmt.Printf("error closing log file: %w", err)
		}
	}()

	if *verbose {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	}

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

	//go func() {
	//	ticker := time.NewTicker(1 * time.Second)
	//	for range ticker.C {
	//		// duration .2s
	//		duration := 200 * time.Millisecond
	//
	//		p.Send(tea.Tick(duration, func(t time.Time) tea.Msg {
	//			print("tick")
	//			return model.OutputMsg{}
	//		}))
	//	}
	//}()

	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running program: %v", err)
	}

	fmt.Println("\nthat was cool i guess, all is stopped now though")
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
