package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ConfirmDialog struct {
	message  string
	selected bool
	visible  bool
}

func NewConfirmDialog(message string) ConfirmDialog {
	return ConfirmDialog{
		message: message,
		visible: true,
	}
}

func (c ConfirmDialog) Update(msg tea.Msg) (ConfirmDialog, tea.Cmd) {
	if !c.visible {
		return c, nil
	}
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left", "h":
			c.selected = true
		case "right", "l":
			c.selected = false
		case "enter":
			c.visible = false
		case "q", "esc":
			c.visible = false
		}
	}
	return c, nil
}

func (c ConfirmDialog) View() string {
	if !c.visible {
		return ""
	}

	var yesStyle, noStyle lipgloss.Style
	if c.selected {
		yesStyle = lipgloss.NewStyle().Background(lipgloss.Color("12"))
		noStyle = lipgloss.NewStyle()
	} else {
		yesStyle = lipgloss.NewStyle()
		noStyle = lipgloss.NewStyle().Background(lipgloss.Color("12"))
	}
	dialog := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(1, 2).
		Render(c.message + "\n\n" + yesStyle.Render("Yes") + " " + noStyle.Render("No"))

	return dialog
}
