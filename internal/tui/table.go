package tui

import (
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/thanhphuocnguyen/git-branch-cleaner/internal/gitcontrol"
)

var (
	tableStyle = lipgloss.NewStyle().
		Margin(1, 2).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240"))
)

type BranchTable struct {
	table     table.Model
	selected  map[int]struct{}
	branches  []gitcontrol.Branch
	cursor    int
	navigator Navigator
}

func NewTableModel(branches []gitcontrol.Branch) BranchTable {
	columns := []table.Column{
		{Title: "✓", Width: 3},
		{Title: "Branch Name", Width: 25},
		{Title: "IsCurrent", Width: 8},
		{Title: "Hash", Width: 12},
	}
	rows := make([]table.Row, len(branches))

	for i, branch := range branches {
		tickText := " "
		if branch.IsCurrent {
			tickText = "-"
		}
		rows[i] = table.Row{
			tickText, // Selection indicator (will be "✓" when selected)
			branch.Name,
			boolToYesNo(branch.IsCurrent),
			branch.Hash[:8], // Show first 8 chars of hash
		}
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(7),
		table.WithWidth(50),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	m := BranchTable{table: t, selected: make(map[int]struct{}), branches: branches, navigator: NewNavigator(len(branches))}
	return m
}

func (m BranchTable) Init() tea.Cmd {
	return nil
}

func (m BranchTable) Update(msg tea.Msg) (BranchTable, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "space":
			cursor := m.table.Cursor()
			if cursor < len(m.branches) && m.branches[cursor].IsCurrent {
				return m, nil
			}
			if _, ok := m.selected[m.table.Cursor()]; ok {
				delete(m.selected, m.table.Cursor())
			} else {
				m.selected[m.table.Cursor()] = struct{}{}
			}
			m.cursor = cursor
			m.updateSelectionDisplay()
			m.navigator = m.navigator.Update(m.SelectedCount())
		case "q", "ctrl+c":
			return m, tea.Quit
		case "?":
			m.navigator = m.navigator.ToggleHelp()
			return m, nil

		case "c":
			// Clear all selections
			m.selected = make(map[int]struct{})
			m.updateSelectionDisplay()
			m.navigator = m.navigator.Update(0)
			return m, nil

		case "d", "enter":
			if m.SelectedCount() > 0 {
				m.navigator = m.navigator.SetMode("confirm")
			}
			return m, nil
		default:
			if m.navigator.showHelp {
				m.navigator = m.navigator.ToggleHelp()
				return m, nil
			}
		}
	case tea.WindowSizeMsg:
		// Update table and navigator dimensions
		m.table.SetWidth(msg.Width - 4)
		m.table.SetHeight(msg.Height - 8)           // Leave space for navigator
		m.navigator = m.navigator.SetMode("browse") // Reset navigator mode on resize

	}
	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m BranchTable) View() tea.View {
	if m.navigator.showHelp {
		return tea.NewView(m.navigator.View())
	}
	tableView := tableStyle.Render(m.table.View())
	navigatorView := m.navigator.View()
	return tea.NewView(
		lipgloss.JoinVertical(lipgloss.Left, tableView, navigatorView),
	)
}

func boolToYesNo(b bool) string {
	if b {
		return "Yes"
	}
	return "No"
}

func (m *BranchTable) updateSelectionDisplay() {
	rows := m.table.Rows()
	for i := range rows {
		if _, selected := m.selected[i]; selected {
			rows[i][0] = "✓"
		} else {
			rows[i][0] = " "
		}
	}

	m.table.SetRows(rows)
	m.table.SetCursor(m.cursor) // Ensure cursor stays in sync with navigator
}

func (m BranchTable) SelectedIndices() []int {
	indices := make([]int, 0, len(m.selected))
	for idx := range m.selected {
		indices = append(indices, idx)
	}
	return indices
}

func (m BranchTable) SelectedCount() int {
	return len(m.selected)
}
