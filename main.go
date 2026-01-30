package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/thanhphuocnguyen/git-branch-cleaner/internal/gitcontrol"
	"github.com/thanhphuocnguyen/git-branch-cleaner/internal/tui"
)

func main() {
	p := tea.NewProgram(initialModel())
	if err := p.Start(); err != nil {
		panic(err)
	}
}

func initialModel() tea.Model {
	repo, err := gitcontrol.OpenRepository(".")
	if err != nil {
		fmt.Println("Error opening repository: %w\n", err)
		os.Exit(1)
	}

	branches, err := repo.ListLocalBranches()
	if err != nil {
		fmt.Println("Error listing branches: %w\n", err)
		os.Exit(1)
	}

	items := make([]list.Item, len(branches))
	for i, br := range branches {
		items[i] = tui.Item{
			Title:       br.Name,
			Description: fmt.Sprintf("Hash: %s", br.Hash),
			Deletable:   !br.IsCurrent,
		}
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
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
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
		var cmd tea.Cmd
		m.confirmdialog, cmd = m.confirmdialog.Update(msg)
		return m, cmd
	case "deleting":
		// No updates in deleting state.
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case "list":
		return m.multiselect.View()
	case "confirm":
		return m.confirmdialog.View()
	case "deleting":
		return "Deleting selected branches..."
	default:
		return "Unknown state"
	}
}
