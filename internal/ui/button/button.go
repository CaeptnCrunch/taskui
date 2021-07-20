package button

import (
	"github.com/CaeptnCrunch/taskui/internal/styles"
	"strings"
)

type Model struct {
	text      string
	focussed  bool
	acticated bool
}

func NewModel(text string) Model {
	return Model{
		text: text,
	}
}

func (m *Model) SetText(text string) {
	m.text = text
}

func (m *Model) Blur() {
	m.focussed = false
}

func (m *Model) Focus() {
	m.focussed = true
}

func (m *Model) Activate() {
	m.acticated = true
}

func (m *Model) View() string {
	var b strings.Builder

	b.WriteString(styles.Prompt)
	b.WriteString(m.text)
	b.WriteString(styles.AntiPrompt)

	if m.focussed {
		return styles.FocusedStyle.Render(b.String())
	} else if m.acticated {
		return styles.NoStyle.Render(b.String())
	} else {
		return styles.BlurredStyle.Render(b.String())
	}
}
