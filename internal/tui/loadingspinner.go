package tui

// loading view
import (
	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type LoadingSpinner struct {
	spinner spinner.Model
	label   string
}

var spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

func NewLoadingSpinner(label string) LoadingSpinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = spinnerStyle
	return LoadingSpinner{spinner: s, label: label}
}

func (l LoadingSpinner) Init() tea.Cmd {
	return l.spinner.Tick
}

func (l LoadingSpinner) Update(msg tea.Msg) (LoadingSpinner, tea.Cmd) {
	var cmd tea.Cmd
	l.spinner, cmd = l.spinner.Update(msg)
	return l, cmd
}

func (l LoadingSpinner) View() tea.View {
	return tea.NewView(l.spinner.View() + " " + l.label)
}
