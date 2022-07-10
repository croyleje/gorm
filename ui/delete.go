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

type deleteModel struct {
	help         help.Model
	confirmInput textinput.Model
}

func newDeleteModel() *deleteModel {
	ci := textinput.New()
	ci.CharLimit = 1

	return &deleteModel{
		help:         help.New(),
		confirmInput: ci,
	}
}

func deleteUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.delete.confirmInput.Value() {
			case "y", "Y":
				if i, ok := m.list.SelectedItem().(item); ok {
					out := (i.Name)
					cmd.DeleteItem(i.Name)

					m.updateListItem()
					m.delete.confirmInput.Reset()
					m.state = browsing
					m.keyMap.State = "browsing"
					m.updateKeybindings()
					cmd := m.list.NewStatusMessage("deleted: " + out)
					return m, cmd
				}

			case "n", "N", "":
				m.delete.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()

			default:
				m.delete.confirmInput.SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.delete.confirmInput.Reset()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.list.ResetFilter()
		}
	case tea.WindowSizeMsg:
		m.delete.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.delete.confirmInput, cmd = m.delete.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) deleteView() string {
	title := m.styles.Title.MarginLeft(2).Render("delete selected items")
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.delete.help.View(m.keyMap))

	var itemName string

	if i, ok := m.list.SelectedItem().(item); ok {
		itemName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)
	}

	label := fmt.Sprintf("Confirm deletion of items \"%s\"? [y/N]", itemName)

	confirmInput := lipgloss.NewStyle().
		MarginLeft(4).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			label,
			m.delete.confirmInput.View(),
		))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", confirmInput, "\n", help))
}
