package tui

import (
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// This component uses BubbleTea's simple list, like Select
// https://github.com/charmbracelet/bubbles#list

// ErrCanceled is returned when the user quits instead of confirming,
// with q, esc or ctrl+c
var ErrCanceled = errors.New("tui: canceled")

type MultiSelectOptions struct {
	Title string
	Items []SelectItem
	// The items that are selected at the start. The user can unselect them.
	// Those that are not in Items are ignored
	Selected []SelectItem
	// If this function provided, the list will be replaced with the
	// return string once the user confirms
	SelectMessage   func([]string) string
	QuitMessage     string
	ShowStatusBar   bool
	EnableFiltering bool
}

var (
	multiSelectToggleKey  = key.NewBinding(key.WithKeys(" "), key.WithHelp("space", "toggle"))
	multiSelectConfirmKey = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "confirm"))
)

// Shows the items with checkboxes, and returns the selected ones, in the order of Items.
//
// The user moves with the arrow keys, selects or unselects an item with space,
// and confirms with enter. Quitting with q, esc or ctrl+c returns ErrCanceled,
// which tells it apart from confirming with nothing selected
func MultiSelect(options MultiSelectOptions) ([]string, error) {
	v, err := tea.NewProgram(newMultiSelectModel(options)).Run()
	if err != nil {
		return nil, err
	}
	m := v.(multiSelectModel)
	if m.canceled {
		return nil, ErrCanceled
	}
	return m.selection(), nil
}

// An item of the list, with its index in Items, which tells the
// items apart when they are filtered or have the same text
type multiSelectItem struct {
	index int
	label string
}

func (i multiSelectItem) FilterValue() string { return i.label }

type multiSelectDelegate struct {
	// Whether each item of Items is selected. It is shared with the model
	selected []bool
}

func (d multiSelectDelegate) Height() int                             { return 1 }
func (d multiSelectDelegate) Spacing() int                            { return 0 }
func (d multiSelectDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d multiSelectDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(multiSelectItem)
	if !ok {
		return
	}

	box := "[ ]"
	if d.selected[i.index] {
		box = "[x]"
	}
	str := box + " " + i.label

	if index == m.Index() {
		fmt.Fprint(w, selectedItemStyle.Render("> "+str))
		return
	}
	fmt.Fprint(w, itemStyle.Render(str))
}

type multiSelectModel struct {
	list     list.Model
	selected []bool
	done     bool
	canceled bool
	options  MultiSelectOptions
}

func newMultiSelectModel(options MultiSelectOptions) multiSelectModel {
	const defaultWidth = 20

	selected := make([]bool, len(options.Items))
	listItems := []list.Item{}
	for i, v := range options.Items {
		selected[i] = slices.Contains(options.Selected, v)
		listItems = append(listItems, multiSelectItem{index: i, label: string(v)})
	}

	l := list.New(listItems, multiSelectDelegate{selected: selected}, defaultWidth, listHeight)
	l.Title = options.Title
	l.SetShowStatusBar(options.ShowStatusBar)
	l.SetFilteringEnabled(options.EnableFiltering)
	l.Filter = containsFilter
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{multiSelectToggleKey, multiSelectConfirmKey}
	}
	l.AdditionalFullHelpKeys = l.AdditionalShortHelpKeys

	return multiSelectModel{list: l, selected: selected, options: options}
}

func (m multiSelectModel) Init() tea.Cmd {
	return nil
}

func (m multiSelectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if key.Matches(msg, m.list.KeyMap.ForceQuit) {
			m.canceled = true
			return m, tea.Quit
		}

		// While the filter is typed, the keys go to the filter
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch {
		case key.Matches(msg, multiSelectToggleKey):
			if i, ok := m.list.SelectedItem().(multiSelectItem); ok {
				m.selected[i.index] = !m.selected[i.index]
			}
			return m, nil

		case key.Matches(msg, multiSelectConfirmKey):
			m.done = true
			return m, tea.Quit

		// The list clears an applied filter with esc, which is also a quit key
		case key.Matches(msg, m.list.KeyMap.ClearFilter):

		case key.Matches(msg, m.list.KeyMap.Quit):
			m.canceled = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m multiSelectModel) View() string {
	if m.done && m.options.SelectMessage != nil {
		return quitTextStyle.Render(m.options.SelectMessage(m.selection()))
	}
	if m.canceled && m.options.QuitMessage != "" {
		return quitTextStyle.Render(m.options.QuitMessage)
	}
	return "\n" + m.list.View()
}

// Returns the selected items, in the order of Items
func (m multiSelectModel) selection() []string {
	var s []string
	for i, v := range m.options.Items {
		if m.selected[i] {
			s = append(s, string(v))
		}
	}
	return s
}

// Keeps the items that contain the term, ignoring the case, like the filter of Select
func containsFilter(term string, targets []string) []list.Rank {
	term = strings.ToLower(term)
	ranks := []list.Rank{}
	for i, t := range targets {
		if strings.Contains(strings.ToLower(t), term) {
			ranks = append(ranks, list.Rank{Index: i, MatchedIndexes: []int{}})
		}
	}
	return ranks
}
