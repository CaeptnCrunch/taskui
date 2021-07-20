package main

import (
	"github.com/CaeptnCrunch/taskui/internal/config"
	"github.com/CaeptnCrunch/taskui/internal/ui/view/taskdetails"
	tea "github.com/charmbracelet/bubbletea"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	// Load configuration
	err := config.Load(config.GetDefaultConfigPath())

	// Configuration does not exist, so we create a default config
	if os.IsNotExist(err) {
		log.Info("Configuration does not exist, creating default config")
		err = config.CreateDefaultConfig()
		if err != nil {
			log.WithError(err).Error("Unable to create default config")
			os.Exit(1)
		}
	}

	initLogger()
}

func initLogger() {
	logfile, err := os.OpenFile(config.Get().LogfilePath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(logfile)
	}

	level, err := log.ParseLevel(config.Get().Loglevel)
	if err != nil {
		log.SetLevel(log.InfoLevel)
		log.WithError(err).Warn("Unable to set loglevel from config. Defaulting to INFO")
	} else {
		log.SetLevel(level)
	}

	log.Info("Logger initalized")
}

func main() {

	if err := tea.NewProgram(taskdetails.NewModel()).Start(); err != nil {
		log.WithError(err).Fatal("could not start program")
		os.Exit(1)
	}
}
