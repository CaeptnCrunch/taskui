package logger

import (
	"github.com/CaeptnCrunch/taskui/internal/config"
	log "github.com/sirupsen/logrus"
	"os"
)

func Init() {
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
