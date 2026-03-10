package main

import (
	"fmt"
	"os"

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

	return model{
		state:   "loading",
		repo:    *repo,
		spinner: tui.NewLoadingSpinner("Loading branches..."),
	}
}

type branchesLoadedMsg struct {
	branches []gitcontrol.Branch
	err      error
}

type model struct {
	state   string // "welcome", "list", "confirm", "deleting"
	spinner tui.LoadingSpinner
	welcome tui.WelcomeScreen
	table   tui.BranchTable

	repo     gitcontrol.Repository
	branches []gitcontrol.Branch
}

func (m model) loadBranchesCmd(repo gitcontrol.Repository) tea.Cmd {
	return func() tea.Msg {
		branches, err := repo.ListLocalBranches()
		return branchesLoadedMsg{branches: branches, err: err}
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.loadBranchesCmd(m.repo), m.spinner.Init())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case "loading":
		switch msg := msg.(type) {
		case branchesLoadedMsg:
			if msg.err != nil {
				fmt.Printf("Error listing branches: %v\n", msg.err)
				return m, tea.Quit
			}

			if len(msg.branches) == 0 {
				fmt.Println("No local branches found.")
				return m, tea.Quit
			}

			m.state = "welcome"
			m.branches = msg.branches
			m.welcome = tui.NewWelComeScreen()
			return m, nil
		case tea.KeyPressMsg:
			if msg.String() == "q" || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		}
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case "welcome":
		switch msg.(type) {
		case tea.KeyPressMsg:
			// Any key press moves to the list view
			m.state = "list"
			m.table = tui.NewTableModel(m.branches)
			return m, nil
		}

	case "list":
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			}
		}
		var cmd tea.Cmd
		m.table, cmd = m.table.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() tea.View {
	switch m.state {
	case "loading":
		return m.spinner.View()
	case "welcome":
		return m.welcome.View()
	case "list":
		return m.table.View()

	default:
		return tea.NewView("Unknown state")
	}
}
