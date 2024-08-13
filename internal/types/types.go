package types

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thejawker/rennen/internal/process"
	"sync"
	"time"
)

type Model struct {
	Processes       []*process.Process
	Commands        []*process.Process
	Tabs            []Tab
	ActiveTab       int
	WindowSize      tea.WindowSizeMsg
	Viewport        *viewport.Model
	Mutex           sync.Mutex
	StartedAt       time.Time
	SelectedCommand int
}

type Tab struct {
	Name         string
	Notification bool
	Status       string
}

type ViewModelProvider interface {
	GetViewModel() Model
	GetActiveProcess() *process.Process
	GetActiveTabName() string
	IsOverview() bool
	GetRunTime() string
	GetCommandByName(name string) *process.Process
}
