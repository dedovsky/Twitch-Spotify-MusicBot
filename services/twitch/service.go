package twitch

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gempir/go-twitch-irc/v4"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"net/url"
	"strings"
	"test/test/domain/models"
	"test/test/infrastructure/httpClient"
	"test/test/infrastructure/spotify"
	"test/test/services/tokenService"
)

type Service struct {
	username     string
	clientID     string
	clientSecret string
	refreshToken string

	channelName string
	bot         *twitch.Client
	spotify     *spotify.Repository
	songs       map[string]string

	tokenModel models.Token
}

func Init(repository *spotify.Repository) *Service {
	s := &Service{
		username:     "",
		clientID:     "",
		clientSecret: "",
		refreshToken: "",

		channelName: viper.GetString("twitchName"),
		spotify:     repository,
		songs:       make(map[string]string),
	}
	s.SetupToken()
	return s
}

func (s *Service) ListenAndServe() {
	s.bot = twitch.NewClient(s.username, fmt.Sprint("oauth:", s.tokenModel.AccessToken))
	client := s.bot

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if strings.HasPrefix(message.Message, "!") && log.GetLevel() == log.DebugLevel {
			log.Debug(message)
		}
		switch message.Message {
		case "!ping":
			s.ping(message)
		case "!song":
			s.song(message)
		case "!queue", "!songlist":
			s.queue(message)
		default:
			if strings.HasPrefix(message.Message, "!sr ") {
				s.sr(message)
			}
		}
	})

	client.Join(s.channelName)
	log.Debug("Join: " + s.channelName)
	err := client.Connect()
	if err != nil {
		log.Error("Ошибка подключения к твитчу:", err)
		log.Info("Если ошибка возникает слишком часто, просьба написать либо создателю бота,\nлибо человеку, который вам бота скинул.")
	}
}

func (s *Service) GetExpiresIn() int {
	return s.tokenModel.ExpiresIn
}

func (s *Service) SetupToken() {

	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)
	client := fasthttp.Client{}

	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetRequestURI("https://id.twitch.tv/oauth2/token")
	data := url.Values{
		"client_id":     {s.clientID},
		"client_secret": {s.clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {s.refreshToken},
	}
	req.SetBodyString(data.Encode())
	err := client.Do(req, resp)
	if err != nil {
		log.Fatalf("ошибка обновления токена бота: %w", err)
	}

	log.Debug("Запрос токена бота: " + string(req.Body()))
	log.Debug("resp.Body токена бота: " + string(resp.Body()))

	var entity models.Token
	err = json.Unmarshal(resp.Body(), &entity)
	if err != nil {
		log.Fatalf("ошибка десериализации: %w", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Debug("Body: " + string(resp.Body()))
		log.Fatalf("Error: received non-200 response code: %s", resp.StatusCode())
	}
	s.tokenModel = entity
	go tokenService.RefreshEveryExpiresIn(s)
}
