package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Chat struct {
	Username  string `yaml:"username"`
	LastMsgID int    `yaml:"last-message-id"`
}

type Config struct {
	Chats          []Chat `yaml:"default-chats"`
	Token          string `yaml:"bot-token"`
	AppID          int    `yaml:"api-id"`
	APIHash        string `yaml:"api-hash"`
	RequestLimit   int    `yaml:"requests-limit"`
	RequestDelay   int    `yaml:"requests-delay"`
	MessagesLimit  int    `yaml:"messages-limit"`
	OutputFileName string `yaml:"out-file-name"`
	AdminID        int    `yaml:"admin-id"`
}

func readConfig(path string) (Config, error) {
	cfg := Config{}

	file, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
