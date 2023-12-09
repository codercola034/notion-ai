package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/codercola034/notion-ai/notion"
)

type ChatStartMsg struct {
	result notion.Result
}
type ChatContinueMsg struct {
	result notion.Result
}
type ErrMsg error

func NewCmd(msg tea.Msg) tea.Cmd {
	return func() tea.Msg { return msg }
}

func ChatContinue(result notion.Result) tea.Cmd {
	return NewCmd(ChatContinueMsg{result})
}

func ChatStart(query string) tea.Cmd {
	result, err := notion.GetCompletion(query)
	if err != nil {
		return NewCmd(ErrMsg(err))
	}

	return NewCmd(ChatStartMsg{result})
}
