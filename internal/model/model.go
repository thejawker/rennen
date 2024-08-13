package model

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/thejawker/rennen/internal/utils"
	"log"
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thejawker/rennen/internal/process"
	"github.com/thejawker/rennen/internal/types"
	"github.com/thejawker/rennen/internal/ui"
)

type Model struct {
	Commands        []*process.Process
	Processes       []*process.Process
	Tabs            []types.Tab
	ActiveTab       int
	WindowSize      tea.WindowSizeMsg
	Viewport        *viewport.Model
	Mutex           sync.Mutex
	StartedAt       time.Time
	SelectedCommand int
}

func New(processes, commands []*process.Process) *Model {
	tabs := make([]types.Tab, len(processes)+1)
	tabs[0] = types.Tab{Name: "overview", Notification: false}
	for i, p := range processes {
		tabs[i+1] = types.Tab{Name: p.Shortname, Notification: false}
	}

	// attach
	return &Model{
		Processes:       processes,
		Commands:        commands,
		SelectedCommand: 0,
		Tabs:            tabs,
		ActiveTab:       0,
		StartedAt:       time.Now(),
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
			return m, m.Shutdown()
		case "tab", "right", "l":
			m.ActiveTab = (m.ActiveTab + 1) % len(m.Tabs)
			return m.ClearNotification(m.ActiveTab)
		case "shift+tab", "left", "h":
			m.ActiveTab = (m.ActiveTab - 1 + len(m.Tabs)) % len(m.Tabs)
			return m.ClearNotification(m.ActiveTab)
		case "up":
			m.SelectedCommand = (m.SelectedCommand - 1 + len(m.Commands)) % len(m.Commands)
		case "down":
			m.SelectedCommand = (m.SelectedCommand + 1) % len(m.Commands)
		case "x":
			if m.ActiveTab > 0 && m.ActiveTab <= len(m.Processes) {
				return m, m.closeProcess(m.GetActiveProcess())
			}
		case "c":
			if m.ActiveTab > 0 && m.ActiveTab <= len(m.Processes) {
				proc := m.GetActiveProcess()
				return m, func() tea.Msg {
					proc.ClearOutput()
					return ProcessUpdateMsg{}
				}
			}
		case "r":
			if m.ActiveTab > 0 && m.ActiveTab <= len(m.Processes) {
				proc := m.GetActiveProcess()
				return m, m.restartProcess(proc)
			}
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

func (m *Model) restartProcess(proc *process.Process) tea.Cmd {
	return func() tea.Msg {
		if err := proc.Restart(); err != nil {
			log.Printf("Error restarting process %s: %v\n", proc.Shortname, err)
		}

		tab := m.GetTabForProcess(proc)
		if tab != nil {
			tab.Status = ""
		}
		return ProcessUpdateMsg{}
	}
}

func (m *Model) closeProcess(proc *process.Process) tea.Cmd {
	return func() tea.Msg {
		if err := proc.Stop(); err != nil {
			log.Printf("Error stopping process %s: %v\n", proc.Shortname, err)
		}
		tab := m.GetTabForProcess(proc)
		if tab != nil {
			tab.Status = "stopped"
			tab.Notification = false
		}
		return ProcessUpdateMsg{}
	}
}

type OutputMsg struct {
	process *process.Process
}

func (m *Model) Shutdown() tea.Cmd {
	return func() tea.Msg {
		var wg sync.WaitGroup
		for _, p := range m.Processes {
			wg.Add(1)
			go func(proc *process.Process) {
				defer wg.Done()
				if err := proc.Stop(); err != nil {
					log.Printf("Error stopping process %s: %v\n", proc.Shortname, err)
				}
			}(p)
		}
		wg.Wait()

		return tea.Quit()
	}
}

func (m *Model) ClearNotification(tabIndex int) (tea.Model, tea.Cmd) {
	var tab types.Tab
	if tabIndex > 0 && tabIndex < len(m.Tabs) {
		tab = m.Tabs[tabIndex]
		tab.Notification = false
	}

	proc := m.GetProcessForTab(tab)

	if proc != nil {
		// Reset the last activity time
		log.Printf("Clearing notification for %s\n", proc.Shortname)
		proc.LastActivity = time.Now().Add(-time.Minute * 10)
	}

	return m, func() tea.Msg {
		return ProcessUpdateMsg{}
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
				if p.Shortname == t.Name {
					if m.GetActiveTabName() == t.Name {
						m.Tabs[i].Notification = false
						break
					}

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

func (m *Model) IsOverview() bool {
	return m.ActiveTab == 0
}

func (m *Model) GetRunTime() string {
	return utils.RelativeTime(m.StartedAt)
}

func (m *Model) startAllProcesses() []tea.Cmd {
	cmds := make([]tea.Cmd, len(m.Processes))
	for i, p := range m.Processes {
		i, p := i, p
		cmds[i] = m.startProcess(p)
	}
	return cmds
}

func (m *Model) startProcess(p *process.Process) func() tea.Msg {
	return func() tea.Msg {
		err := p.Start()
		if err != nil {
			return ProcessErrorMsg{Process: p, Err: err}
		}
		return ProcessStartedMsg{Process: p}
	}
}

func (m *Model) GetTabForProcess(proc *process.Process) *types.Tab {
	for i, t := range m.Tabs {
		if t.Name == proc.Shortname {
			return &m.Tabs[i]
		}
	}
	return nil
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
	time.Sleep(time.Millisecond * 200)
	return tickMsg{}
}
