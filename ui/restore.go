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

type restoreModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newRestoreModel() *restoreModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &restoreModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func restoreUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.restore.confirmInput.Value() {
			case "y", "Y":
				if i, ok := m.list.SelectedItem().(item); ok {
					out := (i.Name)
					cmd.RestoreItem(i.Name)

					m.updateListItem()
					m.restore.confirmInput.Reset()
					m.state = browsing
					m.keyMap.State = "browsing"
					m.updateKeybindings()
					cmd := m.list.NewStatusMessage("restored to original path: " + out)
					return m, cmd
				}

			case "l", "L":
				if i, ok := m.list.SelectedItem().(item); ok {
					out := (i.Name)
					cmd.RestoreItemLocal(i.Name)

					m.updateListItem()
					m.restore.confirmInput.Reset()
					m.state = browsing
					m.keyMap.State = "browsing"
					m.updateKeybindings()
					cmd := m.list.NewStatusMessage("restored locally: " + out)
					return m, cmd
				}

			case "n", "N", "":
				m.restore.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()

			default:
				m.restore.confirmInput.SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.restore.confirmInput.Reset()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.list.ResetFilter()
		}
	case tea.WindowSizeMsg:
		m.restore.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.restore.confirmInput, cmd = m.restore.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) restoreView() string {
	title := m.styles.Title.MarginLeft(2).Render("restore selected items")
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.restore.help.View(m.keyMap))

	var itemName string

	if i, ok := m.list.SelectedItem().(item); ok {
		itemName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)
	}

	label := fmt.Sprintf("Confirm restoration of items \"%s\"? [l/y/N/]", itemName)

	confirmInput := lipgloss.NewStyle().
		MarginLeft(4).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			label,
			m.restore.confirmInput.View(),
		))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", confirmInput, "\n", help))
}
