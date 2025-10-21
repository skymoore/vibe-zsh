package confirm

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestModelInit(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	cmd := m.Init()
	if cmd != nil {
		t.Error("Init should return nil")
	}
}

func TestModelUpdate_YesKey(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	updatedModel, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("Update should return quit command")
	}

	updated := updatedModel.(model)
	if !updated.confirmed {
		t.Error("Model should be confirmed after 'y' key")
	}
	if updated.cancelled {
		t.Error("Model should not be cancelled after 'y' key")
	}
}

func TestModelUpdate_NoKey(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	updatedModel, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("Update should return quit command")
	}

	updated := updatedModel.(model)
	if updated.confirmed {
		t.Error("Model should not be confirmed after 'n' key")
	}
	if !updated.cancelled {
		t.Error("Model should be cancelled after 'n' key")
	}
}

func TestModelUpdate_EnterWithYesSelected(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("Update should return quit command")
	}

	updated := updatedModel.(model)
	if !updated.confirmed {
		t.Error("Model should be confirmed when Enter pressed with Yes selected")
	}
}

func TestModelUpdate_EnterWithNoSelected(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: false,
	}

	msg := tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, cmd := m.Update(msg)

	if cmd == nil {
		t.Error("Update should return quit command")
	}

	updated := updatedModel.(model)
	if updated.confirmed {
		t.Error("Model should not be confirmed when Enter pressed with No selected")
	}
	if !updated.cancelled {
		t.Error("Model should be cancelled when Enter pressed with No selected")
	}
}

func TestModelUpdate_Toggle(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	// Press tab to toggle
	msg := tea.KeyMsg{Type: tea.KeyTab}
	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(model)

	if updated.selected {
		t.Error("Selected should toggle to false after tab")
	}

	// Press tab again to toggle back
	updatedModel, _ = updated.Update(msg)
	updated = updatedModel.(model)

	if !updated.selected {
		t.Error("Selected should toggle back to true after second tab")
	}
}

func TestModelUpdate_ArrowKeys(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: false,
	}

	// Press left arrow to select Yes
	msg := tea.KeyMsg{Type: tea.KeyLeft}
	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(model)

	if !updated.selected {
		t.Error("Left arrow should select Yes")
	}

	// Press right arrow to select No
	msg = tea.KeyMsg{Type: tea.KeyRight}
	updatedModel, _ = updated.Update(msg)
	updated = updatedModel.(model)

	if updated.selected {
		t.Error("Right arrow should select No")
	}
}

func TestModelUpdate_VimKeys(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: false,
	}

	// Press 'h' to select Yes
	msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
	updatedModel, _ := m.Update(msg)
	updated := updatedModel.(model)

	if !updated.selected {
		t.Error("'h' key should select Yes")
	}

	// Press 'l' to select No
	msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
	updatedModel, _ = updated.Update(msg)
	updated = updatedModel.(model)

	if updated.selected {
		t.Error("'l' key should select No")
	}
}

func TestModelView_Confirmed(t *testing.T) {
	m := model{
		command:   "ls -la",
		confirmed: true,
	}

	view := m.View()
	if view != "" {
		t.Error("View should be empty when confirmed")
	}
}

func TestModelView_Cancelled(t *testing.T) {
	m := model{
		command:   "ls -la",
		cancelled: true,
	}

	view := m.View()
	if !strings.Contains(view, "Cancelled") {
		t.Error("View should contain 'Cancelled' when cancelled")
	}
}

func TestModelView_Normal(t *testing.T) {
	m := model{
		command:  "ls -la",
		selected: true,
	}

	view := m.View()

	// Check for key elements
	if !strings.Contains(view, "Execute this command") {
		t.Error("View should contain title")
	}
	if !strings.Contains(view, "ls -la") {
		t.Error("View should contain command")
	}
	if !strings.Contains(view, "Yes") {
		t.Error("View should contain Yes option")
	}
	if !strings.Contains(view, "No") {
		t.Error("View should contain No option")
	}
	if !strings.Contains(view, "ZSH prompt") {
		t.Error("View should contain outcome preview")
	}
}

func TestModelView_OutcomeChanges(t *testing.T) {
	// Test Yes selected
	m := model{
		command:  "ls -la",
		selected: true,
	}
	view := m.View()
	if !strings.Contains(view, "inserted into your ZSH prompt") {
		t.Error("View should show 'inserted' outcome when Yes is selected")
	}

	// Test No selected
	m.selected = false
	view = m.View()
	if !strings.Contains(view, "discarded") {
		t.Error("View should show 'discarded' outcome when No is selected")
	}
}
