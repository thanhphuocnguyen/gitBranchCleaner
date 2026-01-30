package tui

import tea "github.com/charmbracelet/bubbletea"

type KeyMap struct {
	Up     string
	Down   string
	Select string
	Quit   string
	Delete string
	Help   string
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up:     "up",
		Down:   "down",
		Select: " ",
		Quit:   "q",
		Delete: "enter",
		Help:   "?",
	}
}

func HandleKeyPress(msg tea.KeyMsg, keyMap KeyMap) tea.Cmd {
	switch msg.String() {
	case keyMap.Quit, "ctrl+c":
		return tea.Quit
	default:
		return nil
	}
}
