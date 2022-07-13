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

// func deleteUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch {
// 		case key.Matches(msg, m.keyMap.Enter):
// 			switch m.delete.confirmInput.Value() {
// 			case "y", "Y":
// 				if i, ok := m.list.SelectedItem().(item); ok {
// 					out := (i.Name)
// 					cmd.DeleteItem(i.Name)

// 					m.updateListItem()
// 					m.delete.confirmInput.Reset()
// 					m.state = browsing
// 					m.keyMap.State = "browsing"
// 					m.updateKeybindings()
// 					cmd := m.list.NewStatusMessage("deleted: " + out)
// 					return m, cmd
// 				}

// 			case "n", "N", "":
// 				m.delete.confirmInput.Reset()
// 				m.state = browsing
// 				m.keyMap.State = "browsing"
// 				m.updateKeybindings()
// 				m.updateListItem()

// 			default:
// 				m.delete.confirmInput.SetValue("")
// 			}

// 		case key.Matches(msg, m.keyMap.Cancel):
// 			m.delete.confirmInput.Reset()
// 			m.state = browsing
// 			m.keyMap.State = "browsing"
// 			m.updateKeybindings()
// 			m.updateListItem()
// 			m.list.ResetFilter()
// 			cmd := m.list.NewStatusMessage("canceled deletion")
// 			return m, cmd

// 		}
// 	case tea.WindowSizeMsg:
// 		m.delete.help.Width = msg.Width
// 	}

// 	var cmd tea.Cmd
// 	m.delete.confirmInput, cmd = m.delete.confirmInput.Update(msg)

// 	return m, cmd
// }

func deleteUpdate(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Enter):
			switch m.delete.confirmInput.Value() {
			case "y", "Y":
				var n int = 0
				for _, i := range m.list.Items() {
					if i.(item).IsChecked {
						cmd.DeleteItem(i.(item).Name)
						n++
					}
				}

				m.updateListItem()
				m.delete.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()
				cmd := m.list.NewStatusMessage(m.styles.StatusMsg.Render(" Permanently deleted " + fmt.Sprintf("%d", n) + " items... "))
				return m, cmd

			case "n", "N", "":
				m.delete.confirmInput.Reset()
				m.state = browsing
				m.keyMap.State = "browsing"
				m.updateKeybindings()
				m.updateListItem()

			default:
				m.delete.confirmInput.SetValue("")
			}

		case key.Matches(msg, m.keyMap.Cancel):
			m.delete.confirmInput.Reset()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.updateListItem()
			m.list.ResetFilter()
			cmd := m.list.NewStatusMessage(m.styles.StatusMsg.Render(" Deletion Canceled... "))
			return m, cmd

		}
	case tea.WindowSizeMsg:
		m.delete.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.delete.confirmInput, cmd = m.delete.confirmInput.Update(msg)

	return m, cmd
}

// func (m *Model) deleteView() string {
// 	title := m.styles.Title.MarginLeft(2).Render("delete selected items")
// 	help := lipgloss.NewStyle().MarginLeft(4).Render(m.delete.help.View(m.keyMap))

// 	var itemName string

// 	if i, ok := m.list.SelectedItem().(item); ok {
// 		itemName = lipgloss.NewStyle().
// 			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
// 			Render(i.Name)
// 	}

// 	label := fmt.Sprintf("Confirm deletion of items \"%s\"? [y/N]", itemName)

// 	confirmInput := lipgloss.NewStyle().
// 		MarginLeft(4).
// 		Render(lipgloss.JoinHorizontal(
// 			lipgloss.Left,
// 			label,
// 			m.delete.confirmInput.View(),
// 		))

// 	return lipgloss.NewStyle().
// 		MarginTop(1).
// 		Render(lipgloss.JoinVertical(lipgloss.Left, title, "\n", confirmInput, "\n", help))
// }

func (m *Model) deleteView() string {
	title := m.styles.Title.MarginLeft(2).MarginBottom(1).Render("delete selected items")
	help := lipgloss.NewStyle().MarginLeft(4).PaddingTop(2).Render(m.delete.help.View(m.keyMap))

	// var selectedItems []string
	var renderItems string

	// if i, ok := m.list.SelectedItem().(item); ok {
	// 	selectedItems = lipgloss.NewStyle().
	// 		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
	// 		Render(i.Name)
	// }

	for _, i := range m.list.Items() {
		if i.(item).IsChecked {
			renderItems += i.(item).Name + "\n"
		}
	}

	itemsHeader := m.styles.Delete.Render("Selected for permanent deletion:")
	items := m.styles.DeleteItems.Render(renderItems)

	label := fmt.Sprintf("Confirm deletion? [y/N]")

	confirmInput := lipgloss.NewStyle().
		MarginLeft(4).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			label,
			m.delete.confirmInput.View(),
		))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, itemsHeader, items, confirmInput, help))
}
