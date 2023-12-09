package tui

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	messages messages

	err           error
	textarea      textarea.Model
	viewport      viewport.Model
	yankedByteLen int
}

func NewModel() tea.Model {
	return model{
		textarea: createTextarea(viewportWidth, 5),
		viewport: createViewport(viewportWidth, viewportHeight),
		messages: []interface{}{},
		err:      nil,
	}
}

func createViewport(width, height int) viewport.Model {
	vp := viewport.New(width, height)
	return vp
}

func createTextarea(width, height int) textarea.Model {
	ta := textarea.New()
	ta.SetWidth(width)
	ta.SetHeight(height)
	ta.Focus()
	ta.Placeholder = "Ask a question"
	ta.CharLimit = 4000
	return ta
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	if !m.textarea.Focused() {
		m.viewport, vpCmd = m.viewport.Update(msg)
	}

	switch msg := msg.(type) {
	case ErrMsg:
		m.err = msg
	case ChatStartMsg:
		m.messages = append(m.messages, []interface{}{notionChatRoleMsg, ""}...)
		m.viewport.SetContent(m.messages.String())
		m.viewport.GotoBottom()
		return m, ChatContinue(msg.result)
	case ChatContinueMsg:
		out, ended, _ := msg.result.Output()
		if !ended {
			m.messages[len(m.messages)-1] = notionMsg(m.messages.GetString(len(m.messages)) + out)
			m.viewport.SetContent(m.messages.String())
			m.viewport.GotoBottom()
			return m, ChatContinue(msg.result)
		}
		m.textarea.Focus()
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+y":
			content := m.viewport.View()
			lines := strings.Split(content, "\n")
			for i := len(lines) - 1; i >= 0; i-- {
				if strings.HasPrefix(lines[i], "```") {
					lines = lines[0:i]
					break
				}
			}
			for i := len(lines) - 1; i >= 0; i-- {
				if strings.HasPrefix(lines[i], "```") {
					lines = lines[i+1:]
					break
				}
			}
			content = strings.Join(lines, "\n")
			clipboard.WriteAll(content)
			m.yankedByteLen = len(content)
		case "tab":
			if m.textarea.Focused() {
				m.textarea.Blur()
			} else {
				m.textarea.Focus()
			}
			return m, nil
		case "ctrl+p":
			return m, textarea.Paste
		case "ctrl+c", "esc":
			fmt.Println("\033[H\033[2J") // clear screen
			return m, tea.Quit
		case "ctrl+s":
			val := m.textarea.Value()
			if val == "" {
				return m, nil
			}
			m.messages = append(m.messages, userChatRole)
			m.messages = append(m.messages, userMsg(val))

			m.viewport.SetContent(m.messages.String())
			m.viewport.GotoBottom()

			m.textarea.SetValue("")
			m.textarea.Blur()
			return m, ChatStart(val)
		}
	}

	m.viewport.SetContent(m.messages.String())
	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("🚀  Ask Notion AI something") + "\n\n")
	b.WriteString(m.viewport.View() + "\n\n")
	b.WriteString(m.textarea.View() + "\n\n")
	if m.yankedByteLen > 0 {
		b.WriteString(yankedStyle.Render(fmt.Sprintf("Yanked %d bytes to clipboard", m.yankedByteLen)) + "\n")
	} else {
		b.WriteString("\n\n")
	}
	if m.err != nil {
		b.WriteString(errorStyle.Render(m.err.Error()))
	}
	return appStyle.Render(b.String()) + helpStyle.Render("\nc-s send • c-p paste • tab change-focus • c-y yank-last-codeblock • esc/c-c quit")
}
