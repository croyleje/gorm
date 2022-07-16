package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/croyleje/gorm/cmd"
)

type manageModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newManageModel() *manageModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &manageModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func manageUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.manage.confirmInput.Value() {
			case "y", "Y":
				var n int = 0
				for _, i := range m.list.Items() {
					if i.(item).IsChecked {
						cmd.RestoreItem(i.(item).Name)
						n++
					}
				}

				m.updateListItem()
				m.manage.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()
				cmd := m.list.NewStatusMessage(m.styles.StatusMsg.Render(" Selected " + fmt.Sprintf("%d", n) + " items."))
				return m, cmd

			case "n", "N", "":
				m.manage.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()

			default:
				m.manage.confirmInput.SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.manage.confirmInput.Reset()
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
	m.manage.confirmInput, cmd = m.manage.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) manageView() string {
	title := m.styles.Title.MarginLeft(2).Render("Manage Trash")
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.manage.help.View(m.keyMap))

	label := fmt.Sprintf("Confirm selection? [Y/n]")

	confirmInput := lipgloss.NewStyle().
		MarginLeft(4).
		MarginBottom(2).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			label,
			m.manage.confirmInput.View(),
		))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, confirmInput, help))
}
