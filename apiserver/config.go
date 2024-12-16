package apiserver

import (
	"encoding/json"
	"log"
	"os"
)

const (
	defaultKey = `secret-key`
)

type Config struct {
	ServerAddr string
}

func NewConfig() *Config {
	return &Config{
		ServerAddr: "localhost:9090",
	}
}

func (conf *Config) ReadFile(path string) {
	var confInByte []byte

	readConf := &Config{}

	// проверяем файл конфигурации
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		confInByte, err = os.ReadFile(path)
	}
	if err != nil {
		log.Println("Error database config file use default value")
		err = conf.createConfigFile(path)
		if err != nil {
			log.Println("Error: Criate config file:", err.Error())
		}
		return
	}

	json.Unmarshal(confInByte, readConf)

	if readConf.ServerAddr != `` && readConf.ServerAddr != conf.ServerAddr {
		conf.ServerAddr = readConf.ServerAddr
	}
}

func ReadSequreParam(path string) []byte {
	var key []byte
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		key, err = os.ReadFile(path)
	}
	if err != nil {
		log.Println("error read sequrekey file, use default value")
		key = []byte(defaultKey)
		err = os.WriteFile(path, key, 0644)
		if err != nil {
			log.Println("Error: Criate sequre config file:", err.Error())
		}
	}
	return key
}

func (conf *Config) createConfigFile(path string) error {
	confOutByte, err := json.Marshal(conf)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, confOutByte, 0644)
	if err != nil {
		return err
	}
	return nil
}
