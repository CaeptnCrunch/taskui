package button

import (
	"strings"
	styles "taskui/internal"
)

type Model struct {
	text     string
	focussed bool
}

func NewModel(text string) Model {
	return Model{
		text: text,
	}
}

func (m *Model) Blurr() {
	m.focussed = false
}

func (m *Model) Focus() {
	m.focussed = true
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(styles.Prompt)
	b.WriteString(m.text)
	b.WriteString(styles.AntiPrompt)

	if m.focussed {
		return styles.FocusedStyle.Render(b.String())
	} else {
		return styles.BlurredStyle.Render(b.String())
	}
}
