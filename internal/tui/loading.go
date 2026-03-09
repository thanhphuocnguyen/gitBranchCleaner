package tui

// loading view
import (
	tea "charm.land/bubbletea/v2"
)

type LoadingView struct {
	message string
	visible bool
}

func NewLoadingView(message string) LoadingView {
	return LoadingView{
		message: message,
		visible: true,
	}
}

func (l LoadingView) Update(msg tea.Msg) (LoadingView, tea.Cmd) {
	// Currently, loading view does not handle any messages.
	return l, nil
}

func (l LoadingView) View() string {
	if !l.visible {
		return ""
	}
	return l.message + "\n\nLoading..."
}
