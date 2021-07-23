package config

import (
	"encoding/json"
	"fmt"
	"github.com/CaeptnCrunch/taskui/internal/util"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

const (
	defaultConfigDir  = "~/.taskui"
	defaultConfigFile = defaultConfigDir + "/conf.json"
	defaultLogfile    = defaultConfigDir + "/taskui.log"
	defaultTaskrc     = "~/.taskrc"
)

var (
	conf *Config
)

type Config struct {
	TaskrcPath  string
	LogfilePath string
	Loglevel    string
}

func Load() (err error) {
	configPath := GetDefaultConfigPath()

	// check if config file exits
	_, err = os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Create default config
			err = CreateDefaultConfig()
			if err != nil {
				// If we still have an error we let the user handle it
				return err
			}
			// At this point we assume that the config has been created successfully
		}
	}

	contents, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Unable to open config file: %s", configPath))
		err = CreateDefaultConfig()
		if err != nil {

			return
		} else {

		}
	}

	var config Config
	err = json.Unmarshal(contents, &config)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Error unmarshalling config file: %s", configPath))
		return
	}

	conf = &config
	return
}

func Save(configPath string) (err error) {

	serialisedConfig, err := json.Marshal(conf)
	if err != nil {
		log.WithError(err).Error("Unable to serialize config")
		return
	}

	err = ioutil.WriteFile(configPath, serialisedConfig, 0644)
	if err != nil {
		log.WithError(err).Error(fmt.Sprintf("Unable to write config to %s", configPath))
		return
	}

	return
}

func GetDefaultConfigPath() string {
	return util.ExpandUnixPath(defaultConfigFile)
}

func CreateDefaultConfig() (err error) {
	expandedConfigDir := util.ExpandUnixPath(defaultConfigDir)

	fileInfo, err := os.Stat(expandedConfigDir)

	// create config directory if it does not exist yet
	if err != nil && !os.IsExist(err) {
		err = os.Mkdir(expandedConfigDir, 0755)
		if err != nil {
			log.WithError(err).Error(fmt.Sprintf("Unable to create configuration directory under %s", expandedConfigDir))
			return
		}
	}

	// Check if it is a directory
	if !fileInfo.IsDir() {
		log.Error(fmt.Sprintf("%s already exists and is not a directory", expandedConfigDir))
	}

	// Initialize the config
	config := Config{
		TaskrcPath:  defaultTaskrc,
		LogfilePath: util.ExpandUnixPath(defaultLogfile),
		Loglevel:    log.InfoLevel.String(),
	}
	conf = &config
	err = Save(util.ExpandUnixPath(defaultConfigFile))
	if err != nil {
		log.WithError(err).Error("Unable to save default config")
	}

	return
}

func Get() Config {
	return *conf
}
