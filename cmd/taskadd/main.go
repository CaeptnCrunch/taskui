package main

import (
	"fmt"
	"github.com/CaeptnCrunch/taskui/internal/config"
	"github.com/CaeptnCrunch/taskui/internal/logger"
	"github.com/CaeptnCrunch/taskui/internal/ui/createtasks"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Load configuration
	err := config.Load()

	if err != nil {
		// We were not able to load the configuration or create it in case it did not exist
		log.WithError(err).Error(fmt.Sprintf("Config file could not be created or read. Please check %s",
			config.GetDefaultConfigPath()))
		os.Exit(1)
	}

	logger.Init()
}

func main() {

	if err := tea.NewProgram(createtasks.NewModel()).Start(); err != nil {
		log.WithError(err).Fatal("could not start program")
		os.Exit(1)
	}
}
