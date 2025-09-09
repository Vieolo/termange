package termange

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// This component uses BubbleTea's simple list
// https://github.com/charmbracelet/bubbles#list

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type SelectItem string

func (i SelectItem) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(SelectItem)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

type selectModel struct {
	list          list.Model
	choice        string
	quitting      bool
	SelectOptions SelectOptions
}

func (m selectModel) Init() tea.Cmd {
	return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "q":
			m.quitting = true
			return m, tea.Quit

		case "esc":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(SelectItem)
			if ok {
				m.choice = string(i)
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectModel) View() string {
	if m.choice != "" && m.SelectOptions.SelectMessage != nil {
		return quitTextStyle.Render(m.SelectOptions.SelectMessage(string(m.choice)))
	}
	if m.quitting && m.SelectOptions.QuitMessage != "" {
		return quitTextStyle.Render(m.SelectOptions.QuitMessage)
	}
	return "\n" + m.list.View()
}

type SelectOptions struct {
	Title         string
	Items         []SelectItem
	SelectMessage func(string) string
	QuitMessage   string
}

func Select(options SelectOptions) (string, error) {

	const defaultWidth = 20

	listItems := []list.Item{}

	for _, v := range options.Items {
		listItems = append(listItems, SelectItem(v))
	}

	l := list.New(listItems, itemDelegate{}, defaultWidth, listHeight)
	l.Title = options.Title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := selectModel{
		list:          l,
		SelectOptions: options,
	}
	v, err := tea.NewProgram(m).Run()
	return v.(selectModel).choice, err
}
