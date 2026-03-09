package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"github.com/thanhphuocnguyen/git-branch-cleaner/internal/gitcontrol"
	"github.com/thanhphuocnguyen/git-branch-cleaner/internal/tui"
)

// deleteBranchesMsg is a message type to trigger branch deletion
type deleteBranchesMsg struct{}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}

func initialModel() tea.Model {
	repo, err := gitcontrol.OpenRepository(".")
	if err != nil {
		fmt.Printf("Error opening repository: %v\n", err)
		os.Exit(1)
	}

	branches, err := repo.ListLocalBranches()
	if err != nil {
		fmt.Printf("Error listing branches: %v\n", err)
		os.Exit(1)
	}

	if len(branches) == 0 {
		fmt.Println("No local branches found.")
		os.Exit(0)
	}

	items := make([]list.Item, len(branches))

	for i, br := range branches {
		items[i] = tui.NewItem(br.Name, br.Hash, !br.IsCurrent)
	}

	return model{
		repo:        *repo,
		multiselect: tui.NewMultiSelectList(items),
		state:       "list",
		branches:    branches,
	}
}

type model struct {
	state         string // "list", "confirm", "deleting"
	repo          gitcontrol.Repository
	multiselect   tui.MultiSelectList
	confirmdialog tui.ConfirmDialog
	branches      []gitcontrol.Branch
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case "list":
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "shift+enter":
				selected := m.multiselect.SelectedItems()
				if len(selected) > 0 {
					m.confirmdialog = tui.NewConfirmDialog(
						fmt.Sprintf("Delete %d branches?", len(selected)),
					)
					m.state = "confirm"
				}
			case "q", "ctrl+c":
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.multiselect, cmd = m.multiselect.Update(msg)
		return m, cmd
	case "confirm":
		switch msg := msg.(type) {
		case tui.ConfirmResult:
			switch msg.Action {
			case "confirm":
				if msg.Confirmed {
					// User confirmed deletion
					m.state = "deleting"
					return m, tea.Cmd(func() tea.Msg {
						return deleteBranchesMsg{}
					})
				} else {
					// User selected "No", return to list
					m.state = "list"
					return m, nil
				}
			case "cancel":
				// User cancelled (pressed esc/q), return to list
				m.state = "list"
				return m, nil
			}
		}
		var cmd tea.Cmd
		m.confirmdialog, cmd = m.confirmdialog.Update(msg)
		return m, cmd
	case "deleting":
		switch msg.(type) {
		case deleteBranchesMsg:
			selected := m.multiselect.SelectedItems()
			for _, selectedBranch := range selected {
				if err := m.repo.DeleteBranch(selectedBranch.Name); err != nil {
					fmt.Printf("Error deleting branch %s: %v\n", selectedBranch.Name, err)
				}
			}
			m.state = "list"
			m.branches, _ = m.repo.ListLocalBranches() // Refresh branch list after deletion
			items := make([]list.Item, len(m.branches))
			for i, br := range m.branches {
				items[i] = tui.NewItem(br.Name, br.Hash, !br.IsCurrent)
			}
			m.multiselect = tui.NewMultiSelectList(items)
			return m, nil
		}
	}
	return m, nil
}

func (m model) View() tea.View {
	switch m.state {
	case "list":
		return m.multiselect.View()
	case "confirm":
		return m.confirmdialog.View()
	case "deleting":
		return tea.NewView("Deleting selected branches...")
	default:
		return tea.NewView("Unknown state")
	}
}
