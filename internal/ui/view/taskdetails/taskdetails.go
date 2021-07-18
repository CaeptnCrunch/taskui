package taskdetails

import (
	"fmt"
	"github.com/CaeptnCrunch/go-taskwarrior"
	styles "github.com/CaeptnCrunch/taskui/internal"
	"github.com/CaeptnCrunch/taskui/internal/ui/button"
	"github.com/CaeptnCrunch/taskui/internal/ui/choiceitem"
	"github.com/CaeptnCrunch/taskui/internal/ui/singlechoicecombo"
	"github.com/CaeptnCrunch/taskui/internal/ui/textinput"
	teaTextinput "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"strings"
)

const (
	DESC_INDEX       = iota
	PROJ_INDEX       = iota
	PRIO_BTN_INDEX   = iota
	ENTER_BTN_INDEX  = iota
	CANCEL_BTN_INDEX = iota
)

type Model struct {
	inputs             [2]textinput.Model
	enterButton        button.Model
	cancelButton       button.Model
	openPriorityButton button.Model

	// state of the input componentes
	description string
	project     string

	focusIndex int
	tw         *taskwarrior.TaskWarrior

	// subviews
	priorityView       singlechoicecombo.Model
	priorityViewActive bool
	choosenPriority    choiceitem.Choice
	choices            []choiceitem.Choice
}

func NewModel() Model {
	m := Model{
		enterButton:        button.NewModel("ENTER"),
		cancelButton:       button.NewModel("CANCEL"),
		openPriorityButton: button.NewModel("PRIO"),
		focusIndex:         0,
	}

	m.choices = make([]choiceitem.Choice, 3)
	m.choices[0] = choiceitem.Choice{Key: "L", Label: "Low"}
	m.choices[1] = choiceitem.Choice{Key: "M", Label: "Medium"}
	m.choices[2] = choiceitem.Choice{Key: "H", Label: "High"}

	m.priorityView = singlechoicecombo.NewModel(m.choices)
	m.priorityViewActive = false

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

	tw, err := taskwarrior.NewTaskWarrior("~/.taskrc")
	if err != nil {
		log.Fatalln(err)
	}
	m.tw = tw

	return m
}

func (m Model) createTaskwarriorTask() {
	priority := m.choosenPriority.Key

	if priority != "H" && priority != "M" && priority != "L" {
		priority = ""
	}

	task := taskwarrior.Task{
		Description: m.inputs[DESC_INDEX].Value(),
		Project:     m.inputs[PROJ_INDEX].Value(),
		Priority:    priority,
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

func (m Model) Init() tea.Cmd {
	return teaTextinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case singlechoicecombo.ChoiceMsg:
		m.choosenPriority = msg.Choice
		m.openPriorityButton.Activate()
		m.openPriorityButton.SetText(fmt.Sprintf("PRIO:%s", msg.Choice.Key))
	case tea.KeyMsg:

		switch msg.String() {
		case "esc", "ctrl+q":
			return m, tea.Batch(tea.ClearScrollArea, tea.Quit)
		case "enter":

			switch m.focusIndex {
			case ENTER_BTN_INDEX:
				if m.validateDescription() {
					m.createTaskwarriorTask()
				} else {
					break
				}
				return m, tea.Quit
			case PRIO_BTN_INDEX:
				m.description = m.inputs[DESC_INDEX].Value()
				m.project = m.inputs[PROJ_INDEX].Value()
				m.priorityView.SetParent(&m)
				return &m.priorityView, nil
			case CANCEL_BTN_INDEX:
				if m.focusIndex == CANCEL_BTN_INDEX {
					return m, tea.Quit
				}
			}

		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			// select next index
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > CANCEL_BTN_INDEX {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = CANCEL_BTN_INDEX
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

func statusLine() string {
	var b strings.Builder

	b.WriteString("(tab/shift+tab, up/down: select)")
	b.WriteString(styles.Dot)
	b.WriteString("(enter: choose)")
	b.WriteString(styles.Dot)
	b.WriteString("(esc, ctrl+q: quit)")

	return styles.SubtleStyle.Render(b.String())
}

func (m Model) View() string {
	var b strings.Builder

	// show inputs
	b.WriteString(m.inputs[DESC_INDEX].View())
	b.WriteRune('\n')
	b.WriteString(m.inputs[PROJ_INDEX].View())

	m.enterButton.Blur()
	if m.focusIndex == ENTER_BTN_INDEX {
		m.enterButton.Focus()
	}

	m.openPriorityButton.Blur()
	if m.focusIndex == PRIO_BTN_INDEX {
		m.openPriorityButton.Focus()
	}

	m.cancelButton.Blur()
	if m.focusIndex == CANCEL_BTN_INDEX {
		m.cancelButton.Focus()
	}
	b.WriteRune('\n')
	b.WriteString(m.openPriorityButton.View())
	b.WriteRune('\t')
	b.WriteString(m.enterButton.View())
	b.WriteRune('\t')
	b.WriteString(m.cancelButton.View())

	// status line
	b.WriteRune('\n')
	b.WriteString(statusLine())

	return b.String()
}
