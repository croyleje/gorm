package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/croyleje/gorm/cmd"
	"github.com/croyleje/gorm/list"
	"github.com/croyleje/gorm/ui/keys"
	"github.com/croyleje/gorm/ui/styles"
)

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

	styles := styles.DefaultStyles()
	keys := keys.NewKeyMap()

	l := list.New(items, newItemDelegate(keys, &styles), 55, 33)
	l.Title = "Select Trash Directory"
	l.SetStatusBarItemName("Trash directory", "Trash directories")
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
		}
	case tea.WindowSizeMsg:
		m.manage.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.manage.list, cmd = m.manage.list.Update(msg)

	return m, cmd
}

func (m Model) manageView() string {
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.manage.help.View(m.keyMap))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, m.manage.list.View(), help))
}
