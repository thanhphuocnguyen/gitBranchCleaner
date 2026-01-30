package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type Item struct {
	Title       string
	Description string
	Selected    bool
	Deletable   bool
}

func (i Item) GetTitle() string {
	return i.Title
}

func (i Item) GetDescription() string {
	return i.Description
}

func (i Item) FilterValue() string {
	return i.Title
}

type MultiSelectList struct {
	list     list.Model
	selected map[int]bool
}

func NewMultiSelectList(items []list.Item) MultiSelectList {
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select branches to delete (press space to select, d to delete)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	return MultiSelectList{
		list:     l,
		selected: make(map[int]bool),
	}
}

func (m MultiSelectList) Init() tea.Cmd {
	return nil
}

func (m MultiSelectList) Update(msg tea.Msg) (MultiSelectList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case " ":
			idx := m.list.Index()
			if item, ok := m.list.SelectedItem().(Item); ok && item.Deletable {
				m.selected[idx] = !m.selected[idx]
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MultiSelectList) View() string {
	return m.list.View()
}

func (m MultiSelectList) SelectedItems() []int {
	var selected []int
	for idx, isSelected := range m.selected {
		if isSelected {
			selected = append(selected, idx)
		}
	}
	return selected
}
