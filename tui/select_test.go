package tui

import (
	"testing"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func packages() selectModel {
	return newSelectModel(SelectOptions{
		Items:           []SelectItem{"app", "core", "cli", "common"},
		EnableFiltering: true,
	})
}

func TestSelectPicksTheItemUnderTheCursor(t *testing.T) {
	m, cmd := press(packages(), keyDown, keyEnter)

	if m.choice != "core" || !isQuit(cmd) {
		t.Fatalf("choice = %q, quit = %v, want core and a quit", m.choice, isQuit(cmd))
	}
}

func TestSelectQuit(t *testing.T) {
	for name, k := range map[string]tea.KeyMsg{"q": runes("q"), "esc": keyEsc, "ctrl+c": keyCtrlC} {
		m, cmd := press(packages(), k)
		if !m.quitting || m.choice != "" || !isQuit(cmd) {
			t.Errorf("%s: quitting = %v, choice = %q, quit = %v, want a quit without a choice", name, m.quitting, m.choice, isQuit(cmd))
		}
	}
}

func TestSelectKeysGoToTheFilterWhileTyping(t *testing.T) {
	m, _ := press(packages(), runes("/"), runes("q"))

	if m.quitting {
		t.Fatal("typing q in the filter quit")
	}
	if got := m.list.FilterValue(); got != "q" {
		t.Fatalf("filter = %q, want q", got)
	}
}

func TestSelectEnterWhileTypingAcceptsTheFilter(t *testing.T) {
	m, _ := press(packages(), runes("/"), runes("c"), keyEnter)

	if m.quitting || m.choice != "" {
		t.Fatalf("quitting = %v, choice = %q, want the filter accepted", m.quitting, m.choice)
	}
	if m.list.FilterState() != list.FilterApplied {
		t.Fatalf("filter state = %v, want applied", m.list.FilterState())
	}
}

func TestSelectPicksFromTheFilteredItems(t *testing.T) {
	m := packages()
	m.list.SetFilterText("co")

	m, cmd := press(m, keyDown, keyEnter)

	if m.choice != "common" || !isQuit(cmd) {
		t.Fatalf("choice = %q, quit = %v, want common and a quit", m.choice, isQuit(cmd))
	}
}

func TestSelectEscClearsAnAppliedFilter(t *testing.T) {
	m := packages()
	m.list.SetFilterText("co")

	m, _ = press(m, keyEsc)

	if m.quitting {
		t.Fatal("esc quit instead of clearing the filter")
	}
	if m.list.FilterState() != list.Unfiltered {
		t.Fatalf("filter state = %v, want unfiltered", m.list.FilterState())
	}
}
