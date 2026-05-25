package configs

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var AppSettings Configs

func ReadSettings() error {
	fmt.Println("Starting reading settings file")

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("oшибка загрузки .env файла: %w", err)
	}

	configFile, err := os.Open("internal/configs/configs.json")
	if err != nil {
		return fmt.Errorf("Couldn't open config file. Error is: %w", err)
	}

	defer func(configFile *os.File) {
		err = configFile.Close()
		if err != nil {
			log.Fatal("Couldn't close config file. Error is: ", err.Error())
		}
	}(configFile)

	fmt.Println("Starting decoding settings file")
	if err = json.NewDecoder(configFile).Decode(&AppSettings); err != nil {
		return fmt.Errorf("Couldn't decode settings json file. Error is: %w", err)
	}

	return nil
}
