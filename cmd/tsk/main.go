package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"taskui/internal/ui/choiceitem"
	"taskui/internal/ui/singlechoicecombo"
)

func main() {

	var choices = make([]choiceitem.Choice, 3)
	choices[0] = choiceitem.Choice{Key: "L", Label: "LOW"}
	choices[1] = choiceitem.Choice{Key: "M", Label: "MEDIUM"}
	choices[2] = choiceitem.Choice{Key: "H", Label: "HIGH"}

	if err := tea.NewProgram(singlechoicecombo.NewModel(choices)).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
