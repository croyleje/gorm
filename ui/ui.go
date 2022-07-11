package ui

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/croyleje/gorm/cmd"
	"github.com/croyleje/gorm/list"
	"github.com/croyleje/gorm/ui/keys"
	"github.com/croyleje/gorm/ui/styles"
)

const (
	defaultWidth = 55
	listHeight   = 33
)

type item cmd.Item

func (i item) FilterValue() string { return i.Name }

type state int

const (
	browsing state = iota
	deleting
	details
	restoring
)

type Model struct {
	delete  *deleteModel
	detail  *detailModel
	restore *restoreModel
	keyMap  *keys.KeyMap
	list    list.Model
	styles  styles.Styles
	state   state
}

func InitialModel() Model {
	entries, err := cmd.GetEntries()
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	items := []list.Item{}
	for _, e := range entries {
		items = append(items, item{
			Name:         e.Name,
			Type:         e.Type,
			Size:         e.Size,
			Path:         e.Path,
			DeletionDate: e.DeletionDate,
			IsChecked:    e.IsChecked,
		})
	}

	styles := styles.DefaultStyles()
	keys := keys.NewKeyMap()

	l := list.New(items, newItemDelegate(keys, &styles), defaultWidth, listHeight)
	l.Title = "gorm beta v1.3"
	l.StatusMessageLifetime = time.Duration(2) * time.Second
	l.SetShowStatusBar(true)
	l.Styles.PaginationStyle = styles.Pagination
	l.Styles.HelpStyle = styles.Help

	return Model{
		delete:  newDeleteModel(),
		detail:  newDetailModel(),
		restore: newRestoreModel(),
		keyMap:  keys,
		list:    l,
		styles:  styles,
		state:   browsing,
	}

}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) selectedItem() (item, bool) {
	i, ok := m.list.SelectedItem().(item)

	return i, ok
}

func (m *Model) updateListItem() {
	entries, err := cmd.GetEntries()
	if err != nil {
		log.Fatal(err)
	}

	items := []list.Item{}
	for _, e := range entries {
		items = append(items, item{
			Name:         e.Name,
			Type:         e.Type,
			Size:         e.Size,
			Path:         e.Path,
			DeletionDate: e.DeletionDate,
			IsChecked:    e.IsChecked,
		})
	}

	m.list.SetItems(items)
}

func (m *Model) updateKeybindings() {
	if m.list.SettingFilter() {
		m.keyMap.Enter.SetEnabled(false)
	}

	switch m.state {
	case deleting, details, restoring:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Cancel.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		m.keyMap.Quit.SetEnabled(false)
		m.keyMap.Delete.SetEnabled(false)

		m.list.KeyMap.AcceptWhileFiltering.SetEnabled(false)
		m.list.KeyMap.CancelWhileFiltering.SetEnabled(false)
	case browsing:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Delete.SetEnabled(true)
		m.keyMap.Detail.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		// Cancel and browsing state need to be enabled to allow IsChecked
		// resets vaules when leaving another state.
		m.keyMap.Cancel.SetEnabled(true)

	default:
		m.keyMap.Enter.SetEnabled(true)
		m.keyMap.Delete.SetEnabled(true)
		m.keyMap.Detail.SetEnabled(true)
		m.keyMap.ForceQuit.SetEnabled(true)

		m.keyMap.Cancel.SetEnabled(false)
	}

}

func listUpdate(msg tea.Msg, m *Model) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.list.KeyMap.AcceptWhileFiltering):
			m.state = browsing
			m.updateKeybindings()

		case key.Matches(msg, m.keyMap.CursorUp):
			m.list.CursorUp()

		case key.Matches(msg, m.keyMap.CursorDown):
			m.list.CursorDown()

		case key.Matches(msg, m.keyMap.Cancel):
			m.updateListItem()

		case key.Matches(msg, m.keyMap.Select):
			i := m.list.SelectedItem().(item)
			cmd := m.list.SetItem(m.list.Index(), item{Name: i.Name,
				Type:         i.Type,
				Size:         i.Size,
				Path:         i.Path,
				DeletionDate: i.DeletionDate,
				IsChecked:    !i.IsChecked,
			})
			return m, cmd

		case key.Matches(msg, m.keyMap.Delete):
			m.state = deleting
			m.keyMap.State = "deleting"
			m.delete.confirmInput.Focus()
			m.updateKeybindings()

		case key.Matches(msg, m.keyMap.Detail):
			m.state = details
			m.keyMap.State = "details"
			m.updateKeybindings()

		case key.Matches(msg, m.keyMap.Restore):
			m.state = restoring
			m.keyMap.State = "restoring"
			m.restore.confirmInput.Focus()
			m.updateKeybindings()

		case key.Matches(msg, m.keyMap.Enter):
			if i, ok := m.list.SelectedItem().(item); ok {
				cmd := m.list.NewStatusMessage(i.Name)

				return m, cmd
			}

		}
	}

	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.list.SettingFilter() {
		m.keyMap.Enter.SetEnabled(false)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		}
	}
	switch m.state {
	case browsing:
		return listUpdate(msg, &m)

	case deleting:
		return deleteUpdate(msg, m)

	case details:
		return detailUpdate(msg, m)

	case restoring:
		return restoreUpdate(msg, m)

	default:
		return listUpdate(msg, &m)
		// return m, nil
	}

}

func (m Model) View() string {
	switch m.state {
	case browsing:
		return "\n" + m.list.View()

	case deleting:
		return m.deleteView()

	case details:
		return m.detailView()

	case restoring:
		return m.restoreView()

	default:
		return ""
	}
}
