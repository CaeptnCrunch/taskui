package taskdetails

import (
	"fmt"
	"github.com/CaeptnCrunch/go-taskwarrior"
	teaTextinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
	"taskui/internal/ui/button"
	"taskui/internal/ui/textinput"
)

const (
	DESC_INDEX = iota
	PROJ_INDEX = iota
	PRIO_INDEX = iota
)

type Model struct {
	inputs       [3]textinput.Model
	enterButton  button.Model
	cancelButton button.Model
	focusIndex   int
	tw           *taskwarrior.TaskWarrior
}

func NewModel() Model {
	m := Model{
		enterButton:  button.NewModel("ENTER"),
		cancelButton: button.NewModel("CANCEL"),
		focusIndex:   0,
	}

	var t textinput.Model
	// initialize description input
	t = textinput.NewModel()
	t.SetPlaceholder("Description")
	t.Focus()
	m.inputs[DESC_INDEX] = t

	// initialize project input
	t = textinput.NewModel()
	t.SetPlaceholder("Project")
	m.inputs[PROJ_INDEX] = t

	// initialize priority input
	t = textinput.NewModel()
	t.SetPlaceholder("Priority")
	t.SetCharLimit(1)
	m.inputs[PRIO_INDEX] = t

	tw, err := taskwarrior.NewTaskWarrior("~/.taskrc")
	if err != nil {
		log.Fatalln(err)
	}
	m.tw = tw

	return m
}

func (m Model) createTaskwarriorTask() {

	prio := strings.ToUpper(m.inputs[PRIO_INDEX].Value())

	if prio != "H" && prio != "M" && prio != "L" {
		prio = ""
	}

	task := taskwarrior.Task{
		Description: m.inputs[DESC_INDEX].Value(),
		Project:     m.inputs[PROJ_INDEX].Value(),
		Priority:    prio,
	}
	m.tw.AddTask(&task)
	err := m.tw.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func (m Model) validateDescription() bool {
	return len(m.inputs[DESC_INDEX].Value()) > 3
}

func (m Model) validatePriority() bool {
	s := m.inputs[PRIO_INDEX].Value()
	return len(s) == 0 || s == "L" || s == "M" || s == "H"
}

func (m Model) Init() tea.Cmd {
	return teaTextinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "esc", "ctrl+q":
			return m, tea.Quit
		case "tab", "shift+tab", "up", "down", "enter":
			s := msg.String()

			// quit on enter after last line
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if m.validatePriority() && m.validateDescription() {
					m.createTaskwarriorTask()
				} else {
					break
				}
				return m, tea.Quit
			}

			if s == "enter" && m.focusIndex == len(m.inputs)+1 {
				fmt.Println("Exit by cancel")
				return m, tea.Quit
			}

			// select next index
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs)+1 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs) + 1
			}

			// set focus
			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *Model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		_, cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m Model) View() string {
	var b strings.Builder

	// show inputs
	b.WriteString(m.inputs[DESC_INDEX].View())
	b.WriteRune('\n')
	b.WriteString(m.inputs[PROJ_INDEX].View())
	b.WriteString("\t")
	b.WriteString(m.inputs[PRIO_INDEX].View())

	m.enterButton.Blurr()
	if m.focusIndex == len(m.inputs) {
		m.enterButton.Focus()
	}

	m.cancelButton.Blurr()
	if m.focusIndex == len(m.inputs)+1 {
		m.cancelButton.Focus()
	}
	b.WriteRune('\n')
	b.WriteString(m.enterButton.View())
	b.WriteRune('\t')
	b.WriteString(m.cancelButton.View())

	return b.String()
}
