package main

import (
	box2 "github.com/Delta456/box-cli-maker/v2"
	"github.com/charmbracelet/log"
	"os"
	"test/test/infrastructure/spotify"
	"test/test/pkg/viper"
	"test/test/services/twitch"
	"test/test/services/websocket"
	"time"
)

func main() {
	start := time.Now()

	boxPrint()

	if len(os.Args) > 1 {
		log.Debug("Debug mode")
		log.SetLevel(log.DebugLevel)
	}
	viper.Init()

	rep := spotify.NewRepository()

	w := websocket.Init(rep)
	go w.Connect()

	log.Info("Запущен за " + time.Since(start).String())

	tw := twitch.Init(rep)
	tw.ListenAndServe()
}

func boxPrint() {
	box := box2.New(box2.Config{
		Py:           1,
		Px:           10,
		Type:         "Bold",
		ContentAlign: "Center",
		Color:        "Cyan",
		TitlePos:     "Top",
		TitleColor:   "Yellow",
		ContentColor: "Yellow"})
	box.Println("Bot created by", "DEDovsky1\nhttps://t.me/DEDovsky1\nDiscord: dedovsky1")
}
