package main

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
	handlers "natuna.org/copyright/handlers"
)

type Configuration struct {
	Copyright  []string
	Extensions []Extension
}

type Extension struct {
	Extension string
	Processor string
	Protected []string
}

var Config *Configuration

/*
 * Load the configuration file.
 */
func loadConfigurationFile() error {

	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get user home directory: %v\n", err)
		return err
	}

	configPath := home + "/.copyright/config.toml"
	_, err = os.Stat(configPath)
	if err != nil && isVerbose {
		fmt.Fprintf(os.Stdout, "Warning: failed to find global configuration file: %v\n", err)
		return nil
	}

	Config, err = readConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stdout, "Error: failed to read global configuration file: %v\n", err)
		return err
	}

	loadExtensions(Config)
	return nil
}

func readConfig(filename string) (*Configuration, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Configuration
	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func loadExtensions(config *Configuration) {
	for _, ext := range config.Extensions {
		switch ext.Processor {
		case "apt":
			processor.Handlers[ext.Extension] = handlers.AptHandler{}
			FileHandler.AddProtected(handlers.AptHandler{}, ext.Protected)
		case "bat":
			processor.Handlers[ext.Extension] = handlers.BatHandler{}
			FileHandler.AddProtected(handlers.BatHandler{}, ext.Protected)
		case "hashtag":
			processor.Handlers[ext.Extension] = handlers.HashtagHandler{}
			FileHandler.AddProtected(handlers.HashtagHandler{}, ext.Protected)
		case "java":
			processor.Handlers[ext.Extension] = handlers.JavaHandler{}
			FileHandler.AddProtected(handlers.JavaHandler{}, ext.Protected)
		case "xml":
			processor.Handlers[ext.Extension] = handlers.XmlHandler{}
			FileHandler.AddProtected(handlers.XmlHandler{}, ext.Protected)
		default:
			fmt.Fprintf(os.Stderr, "Warning: unknown processor for extension %s\n", ext.Extension)
		}
	}
}
