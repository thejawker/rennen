package process

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/thejawker/rennen/internal/config"
)

// Process represents a running process
type Process struct {
	Shortname    string
	Command      string
	Description  string
	Output       string
	Cmd          *exec.Cmd
	LastActivity time.Time
	mutex        sync.Mutex
}

// InitializeFromConfig creates Process instances from the provided configuration
func InitializeFromConfig(configs []config.ProcessConfig) ([]*Process, error) {
	processes := make([]*Process, len(configs))
	for i, cfg := range configs {
		processes[i] = &Process{
			Shortname:   cfg.Shortname,
			Command:     cfg.Command,
			Description: cfg.Description,
		}
	}
	return processes, nil
}

// Start begins the execution of the process
func (p *Process) Start() error {
	args := strings.Fields(p.Command)
	p.Cmd = exec.Command(args[0], args[1:]...)

	stdout, err := p.Cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	stderr, err := p.Cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	if err := p.Cmd.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	go p.handleOutput(io.MultiReader(stdout, stderr))

	return nil
}

// handleOutput reads the process output and updates the Process struct
func (p *Process) handleOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		p.mutex.Lock()
		p.Output += line + "\n"
		p.LastActivity = time.Now()
		p.mutex.Unlock()
	}
}

// GetOutput returns the current output of the process
func (p *Process) GetOutput() string {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return p.Output
}

// IsActive checks if the process has had activity in the last minute
func (p *Process) IsActive() bool {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	return time.Since(p.LastActivity) < time.Minute
}
