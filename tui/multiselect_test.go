package tui

import (
	"fmt"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	keySpace = tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}
	keyDown  = tea.KeyMsg{Type: tea.KeyDown}
	keyEnter = tea.KeyMsg{Type: tea.KeyEnter}
	keyEsc   = tea.KeyMsg{Type: tea.KeyEsc}
	keyCtrlC = tea.KeyMsg{Type: tea.KeyCtrlC}
)

func runes(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

// Sends the keys to the model, and returns it with the command of the last key
func press(m multiSelectModel, keys ...tea.KeyMsg) (multiSelectModel, tea.Cmd) {
	var cmd tea.Cmd
	for _, k := range keys {
		var next tea.Model
		next, cmd = m.Update(k)
		m = next.(multiSelectModel)
	}
	return m, cmd
}

func isQuit(cmd tea.Cmd) bool {
	if cmd == nil {
		return false
	}
	_, ok := cmd().(tea.QuitMsg)
	return ok
}

func abcd(selected ...SelectItem) multiSelectModel {
	return newMultiSelectModel(MultiSelectOptions{
		Title:    "Pick",
		Items:    []SelectItem{"a", "b", "c", "d"},
		Selected: selected,
	})
}

func TestMultiSelectStartsWithTheInitialSelection(t *testing.T) {
	m := abcd("d", "b", "not an item")

	if got := fmt.Sprint(m.selection()); got != "[b d]" {
		t.Fatalf("selection = %s, want [b d]", got)
	}
}

func TestMultiSelectToggleAndConfirm(t *testing.T) {
	m := abcd("b", "d")

	// Select a, unselect b, then confirm
	m, cmd := press(m, keySpace, keyDown, keySpace, keyEnter)

	if !m.done || m.canceled || !isQuit(cmd) {
		t.Fatalf("done = %v, canceled = %v, quit = %v, want a confirmed quit", m.done, m.canceled, isQuit(cmd))
	}
	if got := fmt.Sprint(m.selection()); got != "[a d]" {
		t.Fatalf("selection = %s, want [a d]", got)
	}
}

func TestMultiSelectKeepsTheOrderOfItems(t *testing.T) {
	m := abcd()

	// Select d, then a
	m, _ = press(m, keyDown, keyDown, keyDown, keySpace)
	m.list.Select(0)
	m, _ = press(m, keySpace)

	if got := fmt.Sprint(m.selection()); got != "[a d]" {
		t.Fatalf("selection = %s, want [a d]", got)
	}
}

func TestMultiSelectConfirmWithNothingSelected(t *testing.T) {
	m, cmd := press(abcd("a"), keySpace, keyEnter)

	if !m.done || m.canceled || !isQuit(cmd) {
		t.Fatalf("done = %v, canceled = %v, want a confirmed quit", m.done, m.canceled)
	}
	if len(m.selection()) != 0 {
		t.Fatalf("selection = %v, want none", m.selection())
	}
}

func TestMultiSelectCancel(t *testing.T) {
	for name, k := range map[string]tea.KeyMsg{"q": runes("q"), "esc": keyEsc, "ctrl+c": keyCtrlC} {
		m, cmd := press(abcd("a"), k)
		if !m.canceled || m.done || !isQuit(cmd) {
			t.Errorf("%s: canceled = %v, done = %v, quit = %v, want a canceled quit", name, m.canceled, m.done, isQuit(cmd))
		}
	}
}

func filterable() multiSelectModel {
	return newMultiSelectModel(MultiSelectOptions{
		Items:           []SelectItem{"app", "core", "cli", "common"},
		EnableFiltering: true,
	})
}

func TestMultiSelectToggleWhileFiltered(t *testing.T) {
	m := filterable()
	m.list.SetFilterText("CO")

	if n := len(m.list.VisibleItems()); n != 2 {
		t.Fatalf("visible items = %d, want core and common", n)
	}
	m, _ = press(m, keyDown, keySpace)

	if got := fmt.Sprint(m.selection()); got != "[common]" {
		t.Fatalf("selection = %s, want [common]", got)
	}
}

func TestMultiSelectEscClearsAnAppliedFilter(t *testing.T) {
	m := filterable()
	m.list.SetFilterText("co")

	m, _ = press(m, keyEsc)

	if m.canceled {
		t.Fatal("esc canceled instead of clearing the filter")
	}
	if m.list.FilterState() != list.Unfiltered {
		t.Fatalf("filter state = %v, want unfiltered", m.list.FilterState())
	}
}

func TestMultiSelectKeysGoToTheFilterWhileTyping(t *testing.T) {
	m := filterable()

	m, _ = press(m, runes("/"), runes("q"), keySpace)

	if m.canceled || m.done {
		t.Fatalf("canceled = %v, done = %v, want the keys typed in the filter", m.canceled, m.done)
	}
	if got := m.list.FilterValue(); got != "q " {
		t.Fatalf("filter = %q, want %q", got, "q ")
	}
	if len(m.selection()) != 0 {
		t.Fatalf("selection = %v, want none", m.selection())
	}
}

func TestMultiSelectView(t *testing.T) {
	view := abcd("b").View()

	for _, want := range []string{"> [ ] a", "[x] b", "[ ] c"} {
		if !strings.Contains(view, want) {
			t.Errorf("view does not contain %q:\n%s", want, view)
		}
	}
}

func TestMultiSelectMessages(t *testing.T) {
	options := MultiSelectOptions{
		Items:         []SelectItem{"a", "b"},
		Selected:      []SelectItem{"b"},
		SelectMessage: func(s []string) string { return "Picked " + strings.Join(s, ", ") },
		QuitMessage:   "Nothing changed",
	}

	m, _ := press(newMultiSelectModel(options), keyEnter)
	if got := m.View(); !strings.Contains(got, "Picked b") {
		t.Errorf("view after confirming = %q, want the select message", got)
	}

	m, _ = press(newMultiSelectModel(options), keyEsc)
	if got := m.View(); !strings.Contains(got, "Nothing changed") {
		t.Errorf("view after canceling = %q, want the quit message", got)
	}
}
