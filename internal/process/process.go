package process

import (
	"fmt"
	"github.com/thejawker/rennen/internal/utils"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"syscall"
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
	done         chan struct{}
	stopped      bool
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

func (p *Process) Restart() error {
	if err := p.Stop(); err != nil {
		return fmt.Errorf("failed to stop process: %w", err)
	}

	p.stopped = false
	p.Output = "restarted process...\n"

	if err := p.Start(); err != nil {
		return fmt.Errorf("failed to start process: %w", err)
	}

	return nil
}

// Start begins the execution of the process
func (p *Process) Start() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.stopped {
		return fmt.Errorf("process has been stopped")
	}

	p.done = make(chan struct{})

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", p.Command)
	} else if shell := os.Getenv("SHELL"); strings.Contains(shell, "zsh") {
		cmd = exec.Command("zsh", "-c", p.Command)
	} else {
		cmd = exec.Command("sh", "-c", p.Command)
	}

	// force color
	cmd.Env = append(os.Environ(),
		"TERM=xterm-256color",
		"COLORTERM=truecolor",
		"FORCE_COLOR=true",
	)

	p.Cmd = cmd

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
	buffer := make([]byte, 1024)
	for {
		select {
		case <-p.done:
			return
		default:
			n, err := reader.Read(buffer)
			if n > 0 {
				p.mutex.Lock()
				p.Output += utils.StripTerminalReturns(string(buffer[:n]))
				p.LastActivity = time.Now()
				p.mutex.Unlock()
			}
			if err != nil {
				if err == io.EOF {
					return
				}
				log.Printf("error reading process output: %v", err)
				return
			}
		}
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

// Stop gracefully stops the process
func (p *Process) Stop() error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.stopped {
		return nil
	}

	p.stopped = true
	close(p.done)

	if p.Cmd != nil && p.Cmd.Process != nil {
		// Send SIGTERM
		err := p.Cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			return fmt.Errorf("failed to send SIGTERM: %w", err)
		}

		p.Output += "\nStopping process\n"

		// Wait for the process to exit or force kill after timeout
		done := make(chan error, 1)
		go func() {
			_, err := p.Cmd.Process.Wait()
			done <- err
		}()

		select {
		case <-time.After(5 * time.Second):
			// Force kill if it doesn't exit within 5 seconds
			err = p.Cmd.Process.Kill()
			if err != nil {
				return fmt.Errorf("failed to kill process: %w", err)
			}
		case err := <-done:
			if err != nil {
				return fmt.Errorf("process exited with error: %w", err)
			}

			p.Output += "\nProcess stopped gracefully\n"
			// Process exited gracefully
		}
	}

	return nil
}

func (p *Process) ClearOutput() {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.Output = ""
}
