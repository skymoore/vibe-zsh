package history

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			Bold(true).
			Foreground(lipgloss.Color("205")). // Bright pink
			Background(lipgloss.Color("235"))  // Dark gray

	itemStyle = lipgloss.NewStyle().
			PaddingLeft(4)

	selectedItemStyle = lipgloss.NewStyle().
				PaddingLeft(2).
				Foreground(lipgloss.Color("86")). // Cyan
				Bold(true)

	timestampStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray

	queryStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")) // Light blue

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")). // Light yellow
			Italic(true)

	arrowStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")) // Gray

	paginationStyle = list.DefaultStyles().PaginationStyle.
			PaddingLeft(4).
			Foreground(lipgloss.Color("241"))

	helpStyle = list.DefaultStyles().HelpStyle.
			PaddingLeft(4).
			PaddingBottom(1).
			Foreground(lipgloss.Color("241"))

	quitTextStyle = lipgloss.NewStyle().
			Margin(1, 0, 2, 4).
			Foreground(lipgloss.Color("203")) // Red
)

type item struct {
	entry Entry
}

type SelectionResult struct {
	Value      string
	Regenerate bool
	EditQuery  bool
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
	timestamp := timestampStyle.Render(i.entry.Timestamp.Format("Jan 02 15:04"))

	// Truncate command if too long
	command := i.entry.Command
	if len(command) > 60 {
		command = command[:57] + "..."
	}

	// Build the item with colors
	query := queryStyle.Render(i.entry.Query)
	arrow := arrowStyle.Render("→")
	cmd := commandStyle.Render(command)

	if index == m.Index() {
		// Selected item - highlight everything
		line1 := fmt.Sprintf("%s  %s", timestamp, query)
		line2 := fmt.Sprintf("    %s %s", arrow, cmd)
		str := fmt.Sprintf("%s\n%s", line1, line2)
		fmt.Fprint(w, selectedItemStyle.Render("▸ "+str))
	} else {
		// Normal item - use individual colors
		line1 := fmt.Sprintf("%s  %s", timestamp, query)
		line2 := fmt.Sprintf("    %s %s", arrow, cmd)
		str := fmt.Sprintf("%s\n%s", line1, line2)
		fmt.Fprint(w, itemStyle.Render(str))
	}
}

type model struct {
	list       list.Model
	choice     string
	regenerate bool
	editQuery  bool
	quitting   bool
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

		case "g":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.entry.Query
				m.regenerate = true
			}
			return m, tea.Quit

		case "v":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = i.entry.Query
				m.editQuery = true
			}
			return m, tea.Quit

		case "a", "home":
			m.list.Select(0)
			return m, nil

		case "e", "end":
			m.list.Select(len(m.list.Items()) - 1)
			return m, nil
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

func ShowInteractive(entries []Entry) (*SelectionResult, error) {
	if len(entries) == 0 {
		return nil, fmt.Errorf("no history entries found")
	}

	// Open /dev/tty directly for TUI interaction
	// This allows the TUI to work even when stdout is captured (e.g., in command substitution)
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// If we can't open /dev/tty, fall back to plain list
		return nil, fmt.Errorf("not in a TTY, use 'vibe-zsh history list' instead")
	}
	defer tty.Close()

	items := make([]list.Item, len(entries))
	for i, entry := range entries {
		items[i] = item{entry: entry}
	}

	const defaultWidth = 80
	const listHeight = 20

	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = " 📜 Query History "
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	// Customize filter prompt colors
	l.Styles.FilterPrompt = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	l.Styles.FilterCursor = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	// Add custom help
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "use command"),
			),
			key.NewBinding(
				key.WithKeys("g"),
				key.WithHelp("g", "regenerate"),
			),
			key.NewBinding(
				key.WithKeys("v"),
				key.WithHelp("v", "edit query"),
			),
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("enter"),
				key.WithHelp("enter", "use generated command"),
			),
			key.NewBinding(
				key.WithKeys("g"),
				key.WithHelp("g", "regenerate from query"),
			),
			key.NewBinding(
				key.WithKeys("v"),
				key.WithHelp("v", "edit query in buffer"),
			),
			key.NewBinding(
				key.WithKeys("a", "home"),
				key.WithHelp("a/home", "go to start"),
			),
			key.NewBinding(
				key.WithKeys("e", "end"),
				key.WithHelp("e/end", "go to end"),
			),
		}
	}

	m := model{list: l}

	// Use /dev/tty for both input and output
	// This keeps stdout free for the selected command output
	// Use AltScreen to avoid messing up the terminal
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty), tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	if m, ok := finalModel.(model); ok {
		if m.choice != "" {
			return &SelectionResult{
				Value:      m.choice,
				Regenerate: m.regenerate,
				EditQuery:  m.editQuery,
			}, nil
		}
	}

	return nil, nil
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
