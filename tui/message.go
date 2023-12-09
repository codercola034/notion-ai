package tui

import "strings"

type chatRoleMsg string
type userMsg string
type notionMsg string

var (
	notionChatRoleMsg = chatRoleMsg("💬 Notion AI")
	userChatRole      = chatRoleMsg("💁 You")
)

type messages []interface{}

func (ms messages) GetString(index int) string {
	switch m := ms[index-1].(type) {
	case notionMsg:
		return string(m)
	case userMsg:
		return string(m) + "\n"
	case chatRoleMsg:
		return string(m) + "\n"
	}
	return ""
}

func (ms messages) Get(index int) interface{} {
	return ms[index-1]
}
func (ms messages) Set(index int, value interface{}) {
	ms[index-1] = value
}

func (ms messages) String() string {
	var builder strings.Builder
	for i := range ms {
		builder.WriteString(ms.GetString(i + 1))
	}
	return builder.String()
}
