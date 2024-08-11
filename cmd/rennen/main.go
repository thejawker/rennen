package main

import (
	"encoding/json"
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

	// handle positional arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "init":
			err := generateDefaultConfig(*configPath)
			if err != nil {
				fmt.Println("Error generating config:", err)
				os.Exit(1)
			}
			fmt.Println("okay, just generated that at", *configPath)
			return
		default:
			fmt.Printf("Unknown command: %s\n", os.Args[1])
			os.Exit(1)
		}
	}

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

func generateDefaultConfig(path string) interface{} {
	// if already exists, panic and exit
	if _, err := os.Stat(path); err == nil {
		log.Fatalf("config file already exists at %s", path)
	} else if !os.IsNotExist(err) {
		log.Fatalf("error checking if config file exists: %v", err)
	}

	defaultConfig := map[string]interface{}{
		// list of objects
		"processes": []map[string]interface{}{
			{
				"shortname":   "test",
				"description": "a sample process",
				"command":     "echo 'hello world'",
			},
		},
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("error closing file: %v", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(defaultConfig)
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
