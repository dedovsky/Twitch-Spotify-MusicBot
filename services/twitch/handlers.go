package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/valyala/fasthttp"
	"strings"
)

func (s *Service) ping(msg twitch.PrivateMessage) {
	s.reply(msg, "pong")
}

func (s *Service) song(msg twitch.PrivateMessage) {
	s.reply(msg, s.spotify.Song())
}

func (s *Service) sr(msg twitch.PrivateMessage) {
	song := strings.TrimPrefix(msg.Message, "!sr ")
	if song == "" {
		s.reply(msg, "Укажите трек")
		return
	}
	if strings.Contains(song, "youtube.com") || strings.Contains(song, "youtu.be") || strings.Contains(song, "open.spotify.com") {
		req := fasthttp.AcquireRequest()
		req.Header.SetMethod(fasthttp.MethodGet)
		if strings.Contains(song, "open.spotify.com") {
			req.SetRequestURI("https://open.spotify.com/oembed?url=" + song)
		} else {
			req.SetRequestURI("https://www.youtube.com/oembed?url=" + song + "&format=json")
		}
		resp := fasthttp.AcquireResponse()
		client := &fasthttp.Client{}
		err := client.Do(req, resp)
		if err != nil {
			s.reply(msg, "Ошибка: "+err.Error())
			log.Error("Ошибка подключения: " + err.Error())
			return
		}
		type responce struct {
			Title string `json:"Title" json:"title"`
		}
		t := responce{Title: ""}
		if err := json.Unmarshal(resp.Body(), &t); err != nil {
			s.reply(msg, "Ошибка: "+err.Error())
			log.Error("Ошибка десериализации: " + err.Error())
			return
		}
		log.Debug("Название песни: " + t.Title)
		song = t.Title
	}

	songModel, err := s.spotify.Sr(song)
	if err != nil {
		s.reply(msg, "Ошибка: "+err.Error())
		log.Error(err)
		return
	} else {
		s.reply(msg, fmt.Sprintf("Трек %s - %s добавлен в очередь", songModel.Artist, songModel.TrackName))
	}
}

func (s *Service) queue(msg twitch.PrivateMessage) {
	s.reply(msg, s.spotify.GetQueue())
}

func (s *Service) reply(msg twitch.PrivateMessage, reply string) {
	text := fmt.Sprintf("@%s, %s", msg.User.DisplayName, reply)
	s.bot.Say(msg.Channel, text)
}
