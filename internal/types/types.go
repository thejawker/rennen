package types

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thejawker/rennen/internal/process"
	"sync"
)

type Model struct {
	Processes  []*process.Process
	Tabs       []Tab
	ActiveTab  int
	WindowSize tea.WindowSizeMsg
	Mutex      sync.Mutex
	Program    *tea.Program
}

type Tab struct {
	Name         string
	Notification bool
}

type ViewModelProvider interface {
	GetViewModel() Model
	GetActiveProcess() *process.Process
	GetActiveTabName() string
}
