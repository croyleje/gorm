package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/croyleje/gorm/cmd"
	"github.com/croyleje/gorm/list"
)

var (
	// titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	// paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	// quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type mountPointDelegate struct{}

func (d mountPointDelegate) Height() int                               { return 1 }
func (d mountPointDelegate) Spacing() int                              { return 0 }
func (d mountPointDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d mountPointDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%v", i.Name)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle.Render("> " + s)
		}
	}

	fmt.Fprintf(w, fn(str))
}

type manageModel struct {
	help help.Model
	list list.Model
}

func newManageModel() *manageModel {
	mounts := cmd.GetMounts()

	items := []list.Item{}
	for _, m := range mounts {
		items = append(items, item{
			Name: m,
		})
	}

	// styles := styles.DefaultStyles()
	// keys := keys.NewKeyMap()

	// l := list.New(items, newItemDelegate(keys, &styles), 55, 33)
	l := list.New(items, mountPointDelegate{}, 55, 33)
	l.Title = "Select Mount Point"
	l.SetStatusBarItemName("Mount Point", "Mount Points")
	l.SetShowStatusBar(true)

	return &manageModel{
		help: help.New(),
		list: l,
	}
}

func manageUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Cancel):
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.updateListItem()
			m.list.ResetFilter()

		case key.Matches(msg, m.keyMap.Enter):
			path := m.manage.list.SelectedItem().(item).Name
			cmd.SetTrashDir(path)
			m.updateListItem()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.list.ResetFilter()
		}
	case tea.WindowSizeMsg:
		m.manage.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.manage.list, cmd = m.manage.list.Update(msg)

	return m, cmd
}

func (m Model) manageView() string {
	help := helpStyle.Render(m.manage.help.View(m.keyMap))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, m.manage.list.View(), help))

}
