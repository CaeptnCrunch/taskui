package singlechoicecombo

import (
	tea "github.com/charmbracelet/bubbletea"
	"strings"
	"taskui/internal/ui/choiceitem"
)

type ChoiceMsg struct {
	Choice choiceitem.Choice
}

type Model struct {
	choices       []choiceitem.Model
	selectedIndex int
	parent        *tea.Model
}

func NewModel(choices []choiceitem.Choice) Model {
	m := Model{
		choices:       make([]choiceitem.Model, len(choices)),
		selectedIndex: 0,
	}

	for i, c := range choices {
		m.choices[i] = choiceitem.NewModel(c)
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) ChoiceMade() tea.Msg {
	return ChoiceMsg{
		Choice: m.choices[m.selectedIndex].GetChoice(),
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			//log.Println("Choice selected")
			return *m.parent, tea.Batch(tea.ClearScrollArea, m.ChoiceMade)

		case "up", "j", "shift+tab":
			m.selectedIndex--
			if m.selectedIndex < 0 {
				m.selectedIndex = len(m.choices) - 1
			}

		case "down", "k", "tab":
			m.selectedIndex = (m.selectedIndex + 1) % len(m.choices)
		}
	}

	return m, nil
}

func (m Model) View() string {

	var b strings.Builder

	for i, c := range m.choices {
		if i == m.selectedIndex {
			c.Focus()
		} else {
			c.Blur()
		}
		b.WriteString(c.View())
		b.WriteRune('\n')
	}

	return b.String()
}

func (m *Model) SetParent(model tea.Model) {
	m.parent = &model
}
