package main

import (
	"fmt"
	"github.com/CaeptnCrunch/go-taskwarrior"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"log"
	"os"
	"strings"
)

const (
	DESC_INDEX = iota
	PROJ_INDEX = iota
	PRIO_INDEX = iota
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	noStyle      = lipgloss.NewStyle()
	prompt       = "❯ "

	enterButtonText      = "❯ ENTER ❮"
	cancelButtonText     = "❯ CANCEL ❮"
	focussedEnterButton  = focusedStyle.Render(enterButtonText)
	blurredEnterButton   = blurredStyle.Render(enterButtonText)
	focussedCancelButton = focusedStyle.Render(cancelButtonText)
	blurredCancelButton  = blurredStyle.Render(cancelButtonText)

	// taskwarrior
	tw *taskwarrior.TaskWarrior
)

type model struct {
	inputs     []textinput.Model
	focusIndex int
}

func initialModel() model {
	m := model{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model

	// initialize description input
	t = textinput.NewModel()
	t.Placeholder = "Description"
	t.Focus()
	t.PromptStyle = focusedStyle
	t.TextStyle = focusedStyle
	t.Prompt = prompt
	m.inputs[DESC_INDEX] = t

	// initialize project input
	t = textinput.NewModel()
	t.Prompt = prompt
	t.Placeholder = "Project"
	m.inputs[PROJ_INDEX] = t

	// initialize priority input
	t = textinput.NewModel()
	t.Prompt = prompt
	t.Placeholder = "Priority"
	t.CharLimit = 1
	m.inputs[PRIO_INDEX] = t

	return m
}

func (m model) createTaskwarriorTask() {

	prio := strings.ToUpper(m.inputs[PRIO_INDEX].Value())

	if prio != "H" && prio != "M" && prio != "L" {
		prio = ""
	}

	task := taskwarrior.Task{
		Description: m.inputs[DESC_INDEX].Value(),
		Project:     m.inputs[PROJ_INDEX].Value(),
		Priority:    prio,
	}
	tw.AddTask(&task)
	err := tw.Commit()
	if err != nil {
		log.Fatalln(err)
	}
}

func (m model) validateDescription() bool {
	return len(m.inputs[DESC_INDEX].Value()) > 3
}

func (m model) validatePriority() bool {
	s := m.inputs[PRIO_INDEX].Value()
	return len(s) == 0 || s == "L" || s == "M" || s == "H"
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
					m.inputs[i].Focus()
					m.inputs[i].TextStyle = focusedStyle
					m.inputs[i].PromptStyle = focusedStyle
				} else {
					m.inputs[i].Blur()
					m.inputs[i].TextStyle = noStyle
					m.inputs[i].PromptStyle = noStyle
				}
			}

			// currently editing the priority
			//val := strings.ToUpper(m.inputs[2].Value())
			//if val != "H" && val != "L" && val != "S" && val != "" {
			//	m.inputs[2].SetValue("")
			//} else {
			//	m.inputs[2].SetValue(strings.ToUpper(val))
			//}

			return m, tea.Batch(cmds...)
		}
	}

	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *model) updateInputs(msg tea.Msg) tea.Cmd {
	var cmds = make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m model) View() string {
	var b strings.Builder

	// show inputs
	b.WriteString(m.inputs[DESC_INDEX].View())
	b.WriteRune('\n')
	b.WriteString(m.inputs[PROJ_INDEX].View())
	b.WriteString("\t")
	b.WriteString(m.inputs[PRIO_INDEX].View())

	//for i := range m.inputs {
	//	b.WriteString(m.inputs[i].View())
	//	if i < len(m.inputs)-1 {
	//		b.WriteRune('\n')
	//	}
	//}

	enterButton := &blurredEnterButton
	if m.focusIndex == len(m.inputs) {
		enterButton = &focussedEnterButton
	}

	cancelButton := &blurredCancelButton
	if m.focusIndex == len(m.inputs)+1 {
		cancelButton = &focussedCancelButton
	}
	b.WriteRune('\n')
	b.WriteString(*enterButton)
	b.WriteRune('\t')
	b.WriteString(*cancelButton)

	return b.String()
}

func main() {

	tskwarrior, err := taskwarrior.NewTaskWarrior("~/.taskrc")
	if err != nil {
		log.Fatalln(err)
	}
	tw = tskwarrior

	if err := tea.NewProgram(initialModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
