package tui

import (
	"fmt"
	"strings"

	"charm.land/lipgloss/v2"
)

type Navigator struct {
	selectedCount int
	totalCount    int
	mode          string
	showHelp      bool
}

func NewNavigator(totalCnt int) Navigator {
	return Navigator{
		selectedCount: 0,
		totalCount:    totalCnt,
		mode:          "browse",
		showHelp:      false,
	}
}

var ( //Add these to your existing styles
	navigatorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Background(lipgloss.Color("235")).
			Padding(0, 1)

	navKeyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Background(lipgloss.Color("235")).
			Bold(true)

	navDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")).
			Background(lipgloss.Color("235"))

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("235")).
			Bold(true)
)

func (n Navigator) Update(selectedCount int) Navigator {
	n.selectedCount = selectedCount
	return n
}

func (n Navigator) ToggleHelp() Navigator {
	n.showHelp = !n.showHelp
	return n
}

func (n Navigator) SetMode(mode string) Navigator {
	n.mode = mode
	return n
}

func (n Navigator) View() string {
	if n.showHelp {
		return n.renderHelp()
	}

	switch n.mode {
	case "browse":
		return n.renderBrowseMode()
	case "confirm":
		return n.renderConfirmMode()
	default:
		return n.renderBrowseMode()
	}
}

func (n Navigator) renderBrowseMode() string {
	var commands []string

	commands = append(commands, keyStyle.Render("space")+" "+descStyle.Render("select"),
		keyStyle.Render("↑↓")+" "+descStyle.Render("navigate"))

	if n.selectedCount > 0 {
		commands = append(commands,
			keyStyle.Render("d")+" "+descStyle.Render("delete"),
			keyStyle.Render("c")+" "+descStyle.Render("clear"),
		)
	}

	commands = append(commands,
		keyStyle.Render("?")+" "+descStyle.Render("help"),
		keyStyle.Render("q")+" "+descStyle.Render("quit"),
	)

	left := lipgloss.JoinHorizontal(lipgloss.Left, commands...)
	// Right side - status
	status := statusStyle.Render(
		fmt.Sprintf("Selected: %d/%d", n.selectedCount, n.totalCount),
	)

	// Calculate spacing to push status to the right
	width := 80 // You can make this dynamic based on terminal width
	spacing := width - lipgloss.Width(left) - lipgloss.Width(status)
	if spacing < 0 {
		spacing = 0
	}

	return navigatorStyle.Width(width).Render(
		left + strings.Repeat(" ", spacing) + status,
	)
}

func (n Navigator) renderConfirmMode() string {
	left := lipgloss.JoinHorizontal(lipgloss.Left,
		keyStyle.Render("y")+" "+descStyle.Render("confirm"),
		keyStyle.Render("n")+" "+descStyle.Render("cancel"),
		keyStyle.Render("esc")+" "+descStyle.Render("back"),
	)

	status := statusStyle.Render(
		fmt.Sprintf("Delete %d branches?", n.selectedCount),
	)

	width := 80
	spacing := width - lipgloss.Width(left) - lipgloss.Width(status)
	if spacing < 0 {
		spacing = 0
	}

	return navigatorStyle.Width(width).Render(
		left + strings.Repeat(" ", spacing) + status,
	)
}

func (n Navigator) renderHelp() string {
	helpText := `
Git Branch Cleaner - Help

Navigation:
  ↑/k        Move up
  ↓/j        Move down
  
Selection:
  space      Toggle branch selection
  c          Clear all selections
  
Actions:
  d          Delete selected branches
  enter      Delete selected branches
  
General:
  ?          Toggle this help
  q/ctrl+c   Quit application
  
Press any key to close help...
`

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("39")).
		Padding(1).
		Margin(1, 2).
		Render(helpText)
}
