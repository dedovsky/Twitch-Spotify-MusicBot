package viper

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

func Init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	if home == "" {
		home = os.Getenv("USERPROFILE")
	}
	viper.AddConfigPath(home + "\\AppData" + "\\Roaming\\tw-mb")
	viper.AddConfigPath(".")
	if viper.ReadInConfig() != nil {
		if runtime.GOOS == "windows" {
			err := os.MkdirAll(home+"\\AppData"+"\\Roaming\\tw-mb", os.ModePerm)
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
			err = viper.SafeWriteConfigAs(home + "\\AppData" + "\\Roaming\\tw-mb\\config.yaml")
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
		} else {
			err := viper.SafeWriteConfigAs("./config.yaml")
			if err != nil {
				log.Fatal("Error: " + err.Error())
			}
		}
	}

	initTwitchName()

	log.Debug("Файл конфига находится в: " + viper.ConfigFileUsed())
}

func initTwitchName() {
	if viper.GetString("twitchName") == "" {
		var userName string
		log.Warn("Не найдено имя аккаунта для твитча в config.yaml. Задайте его")
		_, err := fmt.Scan(&userName)
		if err != nil {
			log.Fatal("Ошибка ввода: " + err.Error())
		}
		viper.Set("twitchName", userName)
		err = viper.WriteConfig()
		if err != nil {
			log.Fatal("Error: " + err.Error())
		}
	} else {
		log.Info("channelName для твитча: " + viper.GetString("twitchName"))
	}
}
