package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	CursorUp   key.Binding
	CursorDown key.Binding
	Select     key.Binding
	SelectAll  key.Binding
	Enter      key.Binding
	Delete     key.Binding
	Detail     key.Binding
	Manage     key.Binding
	Restore    key.Binding
	Help       key.Binding
	Cancel     key.Binding
	Quit       key.Binding
	ForceQuit  key.Binding

	LocalRestore key.Binding

	State string
}

// func (k KeyMap) ShortHelp() []key.Binding {
// 	return []key.Binding{k.Help, k.Quit}
// }
func (k KeyMap) ShortHelp() []key.Binding {
	var kb []key.Binding

	if k.State != "browsing" {
		kb = append(kb, k.Cancel, k.ForceQuit)
	}

	if k.State == "restoring" {
		kb = append(kb, k.LocalRestore)
	}

	return kb
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.CursorUp, k.CursorDown}, // first column
		{k.Help, k.Quit},           // second column
		{k.Cancel, k.ForceQuit},
	}
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("ctrl+k"),
			key.WithHelp("ctrl+k", "move up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("ctrl+j"),
			key.WithHelp("ctrl+j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp(" ", "toggle selection"),
		),
		SelectAll: key.NewBinding(
			key.WithKeys("ctrl+a"),
			key.WithHelp("ctrl+a", "toggle select all"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "confirm selection"),
		),
		Delete: key.NewBinding(
			key.WithKeys("ctrl+d"),
			key.WithHelp("ctrl+d", "delete selected"),
		),
		Detail: key.NewBinding(
			key.WithKeys("ctrl+v"),
			key.WithHelp("ctrl+v", "detail view"),
		),
		Manage: key.NewBinding(
			key.WithKeys("ctrl+t"),
			key.WithHelp("ctrl+t", "choose diffrent trash file"),
		),
		Restore: key.NewBinding(
			key.WithKeys("ctrl+r"),
			key.WithHelp("ctrl+r", "restore selected"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Cancel: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "force quit"),
		),
		LocalRestore: key.NewBinding(
			key.WithKeys("l"),
			key.WithHelp("l", "restore file to current working directory"),
		),
	}
}
