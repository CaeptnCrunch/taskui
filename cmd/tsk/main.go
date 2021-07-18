package main

import (
	"fmt"
	"github.com/CaeptnCrunch/taskui/internal/ui/view/taskdetails"
	tea "github.com/charmbracelet/bubbletea"
	"os"
)

func main() {

	if err := tea.NewProgram(taskdetails.NewModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
