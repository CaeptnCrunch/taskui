package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"taskui/internal/ui/view/taskdetails"
)

func main() {

	if err := tea.NewProgram(taskdetails.NewModel()).Start(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
