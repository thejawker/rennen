package model

import (
	"github.com/charmbracelet/bubbles/viewport"
	"log"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thejawker/rennen/internal/process"
	"github.com/thejawker/rennen/internal/types"
	"github.com/thejawker/rennen/internal/ui"
)

type Model struct {
	Processes  []*process.Process
	Tabs       []types.Tab
	ActiveTab  int
	WindowSize tea.WindowSizeMsg
	Viewport   *viewport.Model
	Mutex      sync.Mutex
}

func New(processes []*process.Process) *Model {
	tabs := make([]types.Tab, len(processes)+1)
	tabs[0] = types.Tab{Name: "Overview", Notification: false}
	for i, p := range processes {
		tabs[i+1] = types.Tab{Name: p.Shortname, Notification: false}
	}

	// attach
	return &Model{
		Processes: processes,
		Tabs:      tabs,
		ActiveTab: 0,
	}
}

func (m *Model) Init() tea.Cmd {
	cmdList := append(m.startAllProcesses(), tick)
	return tea.Batch(cmdList...)
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			m.ClearNotification(m.ActiveTab)
		case "shift+tab", "left", "h":
			m.ActiveTab = (m.ActiveTab - 1 + len(m.Tabs)) % len(m.Tabs)
			m.ClearNotification(m.ActiveTab)
		case "up":
			m.ScrollOutput(-1)
		case "down":
			m.ScrollOutput(1)
		}
	case tea.WindowSizeMsg:
		m.WindowSize = msg
	case ProcessStartedMsg:
		log.Printf("Process started: %s\n", msg.Process.Shortname)
		return m, m.updateNotifications()
	case ProcessErrorMsg:
		log.Printf("Error starting process %s: %v\n", msg.Process.Shortname, msg.Err)
		return m, nil
	case ProcessUpdateMsg, OutputMsg, tickMsg:
		return m, m.updateNotifications()
	}

	return m, nil
}

type OutputMsg struct {
	process *process.Process
}

func (m *Model) ClearNotification(tabIndex int) {
	if tabIndex > 0 && tabIndex < len(m.Tabs) {
		tab := m.Tabs[tabIndex]
		tab.Notification = false
		proc := m.GetProcessForTab(tab)

		if proc != nil {
			// Reset the last activity time
			proc.LastActivity = time.Now().Add(-time.Minute * 2)
		}
	}
}

func (m *Model) GetProcessForTab(tab types.Tab) *process.Process {
	for _, p := range m.Processes {
		if p.Shortname == tab.Name {
			return p
		}
	}
	return nil

}

// updateNotifications checks the last activity time of each process and sets
func (m *Model) updateNotifications() tea.Cmd {
	return func() tea.Msg {
		m.Mutex.Lock()
		defer m.Mutex.Unlock()
		for i, t := range m.Tabs {
			if i == 0 {
				continue // Skip the "all" tab
			}
			for _, p := range m.Processes {
				if p.Shortname == t.Name && m.GetActiveTabName() != t.Name {
					m.Tabs[i].Notification = time.Since(p.LastActivity) < time.Minute
					break
				}
			}
		}

		return tick()
	}
}

func (m *Model) ScrollOutput(amount int) {
	if m.Viewport != nil {
		m.Viewport.LineDown(amount)
	}
}

func (m *Model) View() string {
	return ui.RenderView(m)
}

func (m *Model) GetViewModel() types.Model {
	return types.Model{
		Processes:  m.Processes,
		Tabs:       m.Tabs,
		ActiveTab:  m.ActiveTab,
		WindowSize: m.WindowSize,
	}
}

func (m *Model) GetActiveProcess() *process.Process {
	if m.ActiveTab == 0 || m.ActiveTab > len(m.Processes) {
		return nil
	}
	return m.Processes[m.ActiveTab-1]
}

func (m *Model) GetActiveTabName() string {
	return m.Tabs[m.ActiveTab].Name
}

func (m *Model) startAllProcesses() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Processes))
	for i, p := range m.Processes {
		i, p := i, p // https://golang.org/doc/faq#closures_and_goroutines
		cmds[i] = func() tea.Msg {
			err := p.Start()
			if err != nil {
				return ProcessErrorMsg{Process: p, Err: err}
			}
			return ProcessStartedMsg{Process: p}
		}
	}
	return cmds
}

type ProcessStartedMsg struct {
	Process *process.Process
}

type ProcessErrorMsg struct {
	Process *process.Process
	Err     error
}

type ProcessUpdateMsg struct{}

// Messages are events that we respond to in our Update function. This
// particular one indicates that the timer has ticked.
type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Second)
	return tickMsg{}
}
