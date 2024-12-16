package store

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	DbHost           string `json:"dbhost"`
	DbPort           int    `json:"dbport"`
	DbUser           string `json:"dbuser"`
	DbPassword       string `json:"dbpassword"`
	DbName           string `json:"dbname"`
	DbUsersTable     string `json:"dbutable"`
	DbMessagesTable  string `json:"dbmtable"`
	DbRoomsTable     string `json:"dbrtable"`
	DbNotReadedTable string `json:"dbnrtable"`
}

func NewConfig() *Config {
	return &Config{
		DbHost:           `localhost`,
		DbPort:           5432,
		DbUser:           "postgres",
		DbPassword:       "password",
		DbName:           "web_server",
		DbUsersTable:     "users",
		DbMessagesTable:  "messages",
		DbRoomsTable:     "rooms",
		DbNotReadedTable: "nr_messages",
	}
}

func (conf *Config) ReadConfig(path string) {
	var (
		result       string
		createConfig bool
		confInByte   []byte
	)
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

	// Проверяем пустые поля
	if readConf.DbHost != `` {
		conf.DbHost = readConf.DbHost
	} else {
		result = `default host `
		createConfig = true
	}
	if readConf.DbPort != 0 {
		conf.DbPort = readConf.DbPort
	} else {
		result += `default port `
		createConfig = true
	}
	if readConf.DbUser != `` {
		conf.DbUser = readConf.DbUser
	} else {
		result += `default user `
		createConfig = true
	}
	//пароль не проверяем
	conf.DbPassword = readConf.DbPassword
	if readConf.DbName != `` {
		conf.DbName = readConf.DbName
	} else {
		result += `default name `
		createConfig = true
	}
	if readConf.DbUsersTable != `` {
		conf.DbUsersTable = readConf.DbUsersTable
	} else {
		result += `default users table`
		createConfig = true
	}
	if readConf.DbMessagesTable != `` {
		conf.DbMessagesTable = readConf.DbMessagesTable
	} else {
		result += `default messages table`
		createConfig = true
	}
	if readConf.DbRoomsTable != `` {
		conf.DbRoomsTable = readConf.DbRoomsTable
	} else {
		result += `default messages reference table`
		createConfig = true
	}

	if readConf.DbNotReadedTable != `` {
		conf.DbNotReadedTable = readConf.DbNotReadedTable
	} else {
		result += `default not readed messages table`
		createConfig = true
	}

	if createConfig {
		log.Println("Error database config file use default value:" + result)
		err = conf.createConfigFile(path)
		if err != nil {
			log.Println("Error: Criate config file:", err.Error())
		}
	}
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
