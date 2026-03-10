package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type WelcomeScreen struct {
}

var (
	asciiHeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205"))
	subtitleStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginBottom(1)
	keyStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
	descStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("250"))
	hintStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1).Italic(true)
)

const asciiHeader = `
  +-+-+-+-+ +-+-+-+-+-+-+ +-+-+-+-+-+-+-+
  |G|i|t| |B|r|a|n|c|h| |C|l|e|a|n|e|r|
  +-+-+-+-+ +-+-+-+-+-+-+ +-+-+-+-+-+-+-+
`

func NewWelComeScreen() WelcomeScreen {
	return WelcomeScreen{}
}

func (w WelcomeScreen) Init() tea.Cmd {
	return nil
}

func (w WelcomeScreen) Update(msg tea.Msg) (WelcomeScreen, tea.Cmd) {
	return w, nil
}

func (w WelcomeScreen) View() tea.View {
	content := asciiHeaderStyle.Render(asciiHeader) + "\n" +
		subtitleStyle.Render("  Interactively select and delete local git branches.") + "\n\n" +
		keyStyle.Render("  enter        ") + descStyle.Render("Toggle branch selection") + "\n" +
		keyStyle.Render("  shift+enter  ") + descStyle.Render("Delete selected branches") + "\n" +
		keyStyle.Render("  /            ") + descStyle.Render("Filter branches") + "\n" +
		keyStyle.Render("  q            ") + descStyle.Render("Quit") + "\n" +
		hintStyle.Render("\n  Press any key to continue...")

	return tea.NewView(content)
}
