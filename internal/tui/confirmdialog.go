package tui

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

// ConfirmResult represents the result of a confirm dialog action
type ConfirmResult struct {
	Action    string // "confirm" or "cancel"
	Confirmed bool   // true if user confirmed, false otherwise
}

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
			// Return to list view and trigger delete if selected is true
			if c.selected {
				return c, tea.Cmd(func() tea.Msg {
					return ConfirmResult{Action: "confirm", Confirmed: true}
				})
			} else {
				return c, tea.Cmd(func() tea.Msg {
					return ConfirmResult{Action: "confirm", Confirmed: false}
				})
			}
		case "q", "esc":
			c.visible = false
			return c, tea.Cmd(func() tea.Msg {
				return ConfirmResult{Action: "cancel", Confirmed: false}
			})
		}
	}

	return c, nil
}

var dialogStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2)

func (c ConfirmDialog) View() tea.View {
	if !c.visible {
		return tea.NewView("")
	}

	var yesStyle, noStyle lipgloss.Style
	if c.selected {
		yesStyle = lipgloss.NewStyle().Background(lipgloss.Color("12"))
		noStyle = lipgloss.NewStyle()
	} else {
		yesStyle = lipgloss.NewStyle()
		noStyle = lipgloss.NewStyle().Background(lipgloss.Color("12"))
	}
	dialog :=
		dialogStyle.Render(c.message + "\n\n" + yesStyle.Render("Yes") + " " + noStyle.Render("No"))

	return tea.NewView(dialog)
}
