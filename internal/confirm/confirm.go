package confirm

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")). // Bright pink
			MarginBottom(1)

	commandBoxStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")). // Light yellow
			Background(lipgloss.Color("235")). // Dark gray
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("241")). // Subtle gray border
			MarginBottom(1).
			Width(70) // Fixed width for wrapping

	promptStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("117")) // Light blue

	yesStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("86")) // Cyan

	noStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("203")) // Red

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("205")). // Bright pink
			Bold(true).
			Underline(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")). // Gray
			MarginTop(1)

	outcomeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")). // Subtle gray
			Italic(true).
			MarginTop(1)

	cancelledStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("203")). // Red
			Bold(true)
)

type model struct {
	command   string
	selected  bool // true = yes, false = no
	confirmed bool
	cancelled bool
	width     int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, key.NewBinding(key.WithKeys("ctrl+c", "q", "esc"))):
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("enter"))):
			if m.selected {
				m.confirmed = true
			} else {
				m.cancelled = true
			}
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("y"))):
			m.confirmed = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("n"))):
			m.cancelled = true
			return m, tea.Quit

		case key.Matches(msg, key.NewBinding(key.WithKeys("left", "h"))):
			m.selected = true

		case key.Matches(msg, key.NewBinding(key.WithKeys("right", "l"))):
			m.selected = false

		case key.Matches(msg, key.NewBinding(key.WithKeys("tab"))):
			m.selected = !m.selected
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.confirmed {
		return ""
	}
	if m.cancelled {
		return cancelledStyle.Render("✗ Cancelled") + "\n"
	}

	var s strings.Builder

	// Title with emoji
	s.WriteString(titleStyle.Render("🎯 Execute this command?"))
	s.WriteString("\n\n")

	// Command in a box (will wrap automatically with fixed width)
	s.WriteString(commandBoxStyle.Render(m.command))
	s.WriteString("\n\n")

	// Prompt with Yes/No options
	s.WriteString(promptStyle.Render("Confirm: "))

	var yesText, noText string
	if m.selected {
		yesText = selectedStyle.Render("✓ Yes")
		noText = noStyle.Render("  ✗ No")
	} else {
		yesText = yesStyle.Render("  ✓ Yes")
		noText = selectedStyle.Render("✗ No")
	}

	s.WriteString(yesText)
	s.WriteString("    ")
	s.WriteString(noText)
	s.WriteString("\n")

	// Outcome preview
	if m.selected {
		s.WriteString(outcomeStyle.Render("→ Command will be inserted into your ZSH prompt"))
	} else {
		s.WriteString(outcomeStyle.Render("→ Command will be discarded"))
	}
	s.WriteString("\n")

	// Help
	s.WriteString(helpStyle.Render("← → / h l / tab: toggle • enter: confirm • y: yes • n: no • esc: cancel"))

	return s.String()
}

// ShowConfirmation displays an interactive confirmation prompt for a command
// Returns true if user confirms, false if cancelled
func ShowConfirmation(command string) (bool, error) {
	// Force color output for lipgloss
	lipgloss.SetColorProfile(termenv.TrueColor)

	// Open /dev/tty directly for TUI interaction
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		// Fallback to simple prompt if no TTY
		return simpleConfirm(command)
	}
	defer tty.Close()

	m := model{
		command:  command,
		selected: true, // Default to "Yes" for fast workflow
	}

	// Use /dev/tty for both input and output
	p := tea.NewProgram(m, tea.WithInput(tty), tea.WithOutput(tty))
	finalModel, err := p.Run()
	if err != nil {
		return false, err
	}

	if m, ok := finalModel.(model); ok {
		return m.confirmed && !m.cancelled, nil
	}

	return false, nil
}

// simpleConfirm is a fallback for when TTY is not available
func simpleConfirm(command string) (bool, error) {
	fmt.Fprintf(os.Stderr, "\n---\n%s\n---\n", command)
	fmt.Fprint(os.Stderr, "Execute this command? [Y/n] ")

	var response string
	_, err := fmt.Scanln(&response)
	if err != nil && err != io.EOF {
		return false, err
	}

	// Default to yes if empty or starts with y/Y
	if response == "" || response[0] == 'y' || response[0] == 'Y' {
		return true, nil
	}

	return false, nil
}
