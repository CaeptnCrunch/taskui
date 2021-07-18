package choiceitem

import (
	styles "github.com/CaeptnCrunch/taskui/internal"
	"strings"
)

type Choice struct {
	Key   string
	Label string
}

type Model struct {
	choice  Choice
	focused bool
}

func NewModel(choice Choice) Model {
	return Model{
		choice:  choice,
		focused: false,
	}
}

func (m *Model) Focus() {
	m.focused = true
}

func (m *Model) Blur() {
	m.focused = false
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString(styles.Prompt)
	b.WriteString(m.choice.Label)

	if m.focused {
		return styles.FocusedStyle.Render(b.String())
	} else {
		return styles.NoStyle.Render(b.String())
	}
}

func (m Model) GetChoice() Choice {
	return m.choice
}
