package ui

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/croyleje/gorm/cmd"
)

type detailModel struct {
	help help.Model
}

func newDetailModel() *detailModel {

	return &detailModel{
		help: help.New(),
	}
}

func detailUpdate(msg tea.Msg, m Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Cancel):
			m.delete.confirmInput.Reset()
			m.state = browsing
			m.keyMap.State = "browsing"
			m.updateKeybindings()
			m.updateListItem()
			m.list.ResetFilter()
		}
	case tea.WindowSizeMsg:
		m.delete.help.Width = msg.Width
	}

	var cmd tea.Cmd
	m.delete.confirmInput, cmd = m.delete.confirmInput.Update(msg)

	return m, cmd
}

func (m Model) detailView() string {
	title := m.styles.Title.MarginLeft(2).MarginBottom(2).Render("Detail View")
	help := lipgloss.NewStyle().MarginLeft(4).Render(m.detail.help.View(m.keyMap))

	var itemName string

	if i, ok := m.list.SelectedItem().(item); ok {
		itemName = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"}).
			Render(i.Name)

	}

	header := fmt.Sprintf("Details for \"%s\".", itemName)

	i := m.list.SelectedItem().(item)

	var placeHolder []byte
	var dirList string

	// if i.Type == "Directory" {
	// 	dirList += "Directory contents:\n\n"
	// 	// TODO: rewrite with native go package
	// 	// placeHolder, _ = exec.Command("ls", "-1A", cmd.GetTrashDir()+"files/"+i.Name).Output()
	// 	placeHolder, _ = exec.Command("tree", cmd.GetTrashDir()+"files/"+i.Name).Output()
	// 	dirList += string(placeHolder)
	// 	dirList = m.styles.DetailsFooter.Render(dirList)
	// }

	if i.Type == "Directory" {
		dirList += "Directory contents:\n"
		// TODO: rewrite with native go package
		// placeHolder, _ = exec.Command("ls", "-1A", cmd.GetTrashDir()+"files/"+i.Name).Output()
		placeHolder, _ = exec.Command("tree", cmd.GetTrashDir()+"files/"+i.Name).Output()
		dirList += string(placeHolder)
		dirList = strings.ReplaceAll(dirList, cmd.GetTrashDir()+"files/"+i.Name, "")
		dirList = m.styles.DetailsFooter.Render(dirList + "\n")
	}

	var file string = cmd.GetTrashDir() + "files/" + m.list.SelectedItem().(item).Name
	path := "Path: " + m.list.SelectedItem().(item).Path
	path = m.styles.DetailsHeader.Render(path)

	stats, _ := exec.Command("stat", file).Output()
	statsString := string(stats)
	var formattedString string
	_, formattedString, _ = strings.Cut(statsString, "\n")
	formattedStats := m.styles.Details.Render(string(formattedString))

	render := lipgloss.NewStyle().
		MarginLeft(4).
		MarginBottom(2).
		Render(lipgloss.JoinHorizontal(
			lipgloss.Left,
			header,
		))

	return lipgloss.NewStyle().
		MarginTop(1).
		Render(lipgloss.JoinVertical(lipgloss.Left, title, render, path, formattedStats, dirList, help))
}
