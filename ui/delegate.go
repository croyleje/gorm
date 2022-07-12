package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/croyleje/gorm/list"
	"github.com/croyleje/gorm/ui/keys"
	"github.com/croyleje/gorm/ui/styles"
)

type itemDelegate struct {
	keys   *keys.KeyMap
	styles *styles.Styles
}

func newItemDelegate(keys *keys.KeyMap, styles *styles.Styles) *itemDelegate {
	return &itemDelegate{
		keys:   keys,
		styles: styles,
	}
}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 1 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	title := d.styles.NormalTitle.Render
	desc := d.styles.NormalDesc.Render

	if index == m.Index() {
		title = func(s string) string {
			return d.styles.CurrentTitle.Render("> " + s)
		}
	}

	if i.IsChecked {
		title = func(s string) string {
			return d.styles.SelectedTitle.Render(s)
		}
		desc = func(s string) string {
			return d.styles.SelectedDesc.Render(s)
		}
	}

	if i.IsChecked && index == m.Index() {
		title = func(s string) string {
			return d.styles.SelectedCurrentTitle.Render("> " + s)
		}
		desc = func(s string) string {
			return d.styles.SelectedCurrentDesc.Render(s)
		}
	}

	name := title(i.Name)
	fileType := desc(i.Type)
	deletionDate := desc(i.DeletionDate)
	size := desc(i.Size)
	path := desc(i.Path)

	itemListStyle := fmt.Sprintf("%s %s %s %s %s", name, fileType, size, path, deletionDate)

	fmt.Fprint(w, itemListStyle)
}

func (d itemDelegate) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (d itemDelegate) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{d.keys.Delete, d.keys.Detail, d.keys.Restore, d.keys.SelectAll},
	}
}
