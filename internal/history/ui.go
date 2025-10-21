package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	entry Entry
}

func (i item) FilterValue() string { return i.entry.Query }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 2 }
func (d itemDelegate) Spacing() int                            { return 1 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	// Format timestamp
	timestamp := i.entry.Timestamp.Format("Jan 02 15:04")

	// Truncate command if too long
	command := i.entry.Command
	if len(command) > 60 {
		command = command[:57] + "..."
	}

	str := fmt.Sprintf("%s  %s\n    → %s", timestamp, i.entry.Query, command)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("▸ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type model struct {
	list     list.Model
	choice   string
	quitting bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.entry.Command
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.choice != "" {
		return ""
	}
	if m.quitting {
		return quitTextStyle.Render("Cancelled.")
	}
	return "\n" + m.list.View()
}

func ShowInteractive(entries []Entry) (string, error) {
	if len(entries) == 0 {
		return "", fmt.Errorf("no history entries found")
	}

	// Open /dev/tty directly for TUI interaction
	// This allows the TUI to work even when stdout is captured (e.g., in command substitution)
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// If we can't open /dev/tty, fall back to plain list
		return "", fmt.Errorf("not in a TTY, use 'vibe-zsh history list' instead")
	}
	defer tty.Close()

	items := make([]list.Item, len(entries))
	for i, entry := range entries {
		items[i] = item{entry: entry}
	}

	const defaultWidth = 80
	const listHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Query History"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := model{list: l}

	// Use /dev/tty for both input and output
	// This keeps stdout free for the selected command output
	// Use AltScreen to avoid messing up the terminal
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(model); ok {
		return m.choice, nil
	}

	return "", nil
}

func FormatPlainList(entries []Entry) string {
	if len(entries) == 0 {
		return "No history entries found."
	}

	var sb strings.Builder
	for i, entry := range entries {
		timestamp := entry.Timestamp.Format(time.RFC3339)
		sb.WriteString(fmt.Sprintf("%d. [%s] %s\n", i+1, timestamp, entry.Query))
		sb.WriteString(fmt.Sprintf("   → %s\n\n", entry.Command))
	}
	return sb.String()
}
