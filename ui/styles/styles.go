package styles

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Styles struct {
	Title lipgloss.Style

	NormalTitle lipgloss.Style
	NormalDesc  lipgloss.Style

	SelectedTitle lipgloss.Style
	SelectedDesc  lipgloss.Style

	CurrentTitle lipgloss.Style
	CurrentDesc  lipgloss.Style

	SelectedCurrentTitle lipgloss.Style
	SelectedCurrentDesc  lipgloss.Style

	Details       lipgloss.Style
	DetailsHeader lipgloss.Style
	DetailsFooter lipgloss.Style

	Delete      lipgloss.Style
	DeleteItems lipgloss.Style

	ErrMsg    lipgloss.Style
	StatusMsg lipgloss.Style

	Pagination lipgloss.Style
	Help       lipgloss.Style
	QuitText   lipgloss.Style
}

func DefaultStyles() (s Styles) {
	// Header
	s.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("004")).
		Foreground(lipgloss.Color("007")).
		Padding(0, 1)

	// Title & Description
	s.NormalTitle = lipgloss.NewStyle().
		PaddingLeft(6).
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	s.NormalDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	// Cursor position Title & Description
	s.CurrentTitle = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.AdaptiveColor{Light: "#1A1A1A", Dark: "#DDDDDD"})

	s.CurrentDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#A49FA5", Dark: "#777777"})

	// Cursor position && Selected Title & Description
	s.SelectedCurrentTitle = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	s.SelectedCurrentDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	// Selected Title & Description
	s.SelectedTitle = lipgloss.NewStyle().
		PaddingLeft(6).
		Foreground(lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#EE6FF8"})

	s.SelectedDesc = lipgloss.NewStyle().
		Foreground(lipgloss.AdaptiveColor{Light: "#F793FF", Dark: "#AD58B4"})

	// detailsView
	s.Details = lipgloss.NewStyle().
		PaddingLeft(4).
		Foreground(lipgloss.AdaptiveColor{Light: "#AA6FF8", Dark: "#AA6FF8"})

	s.DetailsHeader = lipgloss.NewStyle().
		PaddingLeft(6).
		Foreground(lipgloss.AdaptiveColor{Light: "#AA6FF8", Dark: "#AA6FF8"})

	s.DetailsFooter = lipgloss.NewStyle().
		PaddingLeft(6).
		Foreground(lipgloss.AdaptiveColor{Light: "#AA6FF8", Dark: "#AA6FF8"})

	// deleteView
	s.Delete = lipgloss.NewStyle().
		PaddingLeft(4).
		PaddingTop(2).
		PaddingBottom(2)

	s.DeleteItems = lipgloss.NewStyle().
		PaddingLeft(4).
		PaddingBottom(1).
		Foreground(lipgloss.AdaptiveColor{Light: "#AA6FF8", Dark: "#AA6FF8"})

	// Error message
	s.ErrMsg = lipgloss.NewStyle().
		Background(lipgloss.Color("001")).
		Foreground(lipgloss.Color("007"))

	// Status message
	s.StatusMsg = lipgloss.NewStyle().
		Background(lipgloss.Color("034")).
		Foreground(lipgloss.Color("007"))

	// Footer
	s.Pagination = list.DefaultStyles().
		PaginationStyle.
		PaddingLeft(4)

	s.Help = list.DefaultStyles().
		HelpStyle.
		PaddingLeft(4).
		PaddingBottom(1)

	s.QuitText = lipgloss.NewStyle().
		Margin(1, 0, 2, 4)

	return s
}
