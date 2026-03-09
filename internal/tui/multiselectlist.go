package tui

import (
	"fmt"

	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type BranchItem struct {
	title     string
	desc      string
	selected  bool
	deletable bool
}

func (i BranchItem) Title() string {
	checked := " "
	if i.selected {
		checked = "x"
	}
	if !i.deletable {
		return fmt.Sprintf("[-] %s (current branch)", i.title)
	}
	return fmt.Sprintf("[%s] %s", checked, i.title)
}

func (i BranchItem) Description() string {
	return i.desc
}

func (i BranchItem) FilterValue() string {
	return i.title
}

type MultiSelectList struct {
	list list.Model
}

func NewItem(title, description string, deletable bool) BranchItem {
	return BranchItem{
		title:     title,
		desc:      description,
		deletable: deletable,
	}
}

func NewMultiSelectList(items []list.Item) MultiSelectList {
	items = append(items, BranchItem{title: "Test", desc: "Tests"})
	width, height := 80, 100
	d := list.NewDefaultDelegate()
	l := list.New(items, d, width, height)
	l.Title = "Select branches to delete (press space to select, d to delete)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.SetDelegate(d)

	return MultiSelectList{
		list: l,
	}
}

func (m MultiSelectList) Init() tea.Cmd {
	return nil
}

func (m MultiSelectList) Update(msg tea.Msg) (MultiSelectList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			idx := m.list.Index()
			if item, ok := m.list.SelectedItem().(BranchItem); ok && item.deletable {
				item.selected = !item.selected
				m.list.SetItem(idx, item)
				// m.selected[idx] = !m.selected[idx]
			}
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MultiSelectList) View() tea.View {
	return tea.NewView(docStyle.Render(m.list.View()))

}

type SelectedBranch struct {
	Index int
	Name  string
}

func (m MultiSelectList) SelectedItems() []SelectedBranch {
	var selected []SelectedBranch
	for idx, item := range m.list.Items() {
		if brItem, ok := item.(BranchItem); ok && brItem.selected {
			selected = append(selected, SelectedBranch{Index: idx, Name: brItem.title})
		}
	}
	return selected
}
