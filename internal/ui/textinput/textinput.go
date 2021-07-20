package textinput

import (
	"github.com/CaeptnCrunch/taskui/internal/styles"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	model *textinput.Model
}

func NewModel() Model {

	m := textinput.NewModel()
	m.Prompt = styles.Prompt
	m.PromptStyle = styles.BlurredStyle
	m.TextStyle = styles.BlurredStyle

	return Model{model: &m}
}

func (m *Model) Blur() {
	m.model.PromptStyle = styles.NoStyle
	m.model.TextStyle = styles.NoStyle
	m.model.Blur()
}

func (m *Model) Focus() tea.Cmd {
	m.model.PromptStyle = styles.FocusedStyle
	m.model.TextStyle = styles.FocusedStyle
	return m.model.Focus()
}

func (m *Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	model, cmd := m.model.Update(msg)
	m.model = &model
	return *m, cmd
}

func (m *Model) SetPlaceholder(str string) {
	m.model.Placeholder = str
}

func (m *Model) SetCharLimit(limit int) {
	m.model.CharLimit = limit
}

func (m Model) Value() string {
	return m.model.Value()
}

func (m Model) View() string {
	return m.model.View()
}

func (m *Model) SetValue(val string) {
	m.model.SetValue(val)
}
