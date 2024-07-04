package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"test/test/domain/models"
	"test/test/infrastructure/httpClient"
	"test/test/infrastructure/spotify"
)

type WebSocket struct {
	channelName string
	tokenModel  models.Token

	url           string
	broadcasterID string
	conn          *websocket.Conn
	httpClient    fasthttp.Client
	clientID      string
	clientSecret  string

	*spotify.Repository
}

func Init(repository *spotify.Repository) *WebSocket {
	w := &WebSocket{
		url:          "wss://eventsub.wss.twitch.tv/ws?keepalive_timeout_seconds=30",
		httpClient:   fasthttp.Client{},
		clientID:     "",
		clientSecret: "",
		channelName:  viper.GetString("twitchName"),

		Repository: repository,
	}
	if viper.GetString("ChannelToken") == "" {
		fmt.Println("Есть ли у вас награда \"Скип трека\"? (y/n)")
		var answer string
		_, err := fmt.Scan(&answer)
		if err != nil {
			log.Fatal("ошибка ввода: " + err.Error())
		}
		if answer == "y" {
			w.SetupToken()
			err = w.initConfig()
			if err != nil {
				log.Fatalf("ошибка подключения WebSocket: %s", err)
			}
		} else {
			viper.Set("channelToken", "nil")
			err = viper.WriteConfig()
			if err != nil {
				log.Fatal("ошибка редактирования конфига: " + err.Error())
			}
			log.Info("Если захотите изменить токен, отредактируйте конфиг.")
		}
	} else if viper.GetString("channelToken") == "nil" {
		return nil
	} else {
		w.SetupToken()
		err := w.initConfig()
		if err != nil {
			log.Fatalf("Ошибка инициализации WebSocket: %v", err)
		}
		w.broadcasterID = viper.GetString("twitchID")
	}

	return w
}

func (w *WebSocket) initConfig() error {
	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://api.twitch.tv/helix/users?login=" + w.channelName)
	req.Header.Set("Client-Id", w.clientID)
	req.Header.Set("Authorization", "Bearer "+w.tokenModel.AccessToken)

	if err := w.httpClient.Do(req, resp); err != nil {
		log.Fatalf("Error sending request: %v", err)
	}

	if resp.StatusCode() != fasthttp.StatusOK {
		log.Debug("token: " + w.tokenModel.AccessToken)
		log.Debug("Body: " + string(resp.Body()))
		log.Fatalf("Error: received non-200 response code: %v", resp.StatusCode())
	}

	var entity TwitchChannelID
	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		log.Fatalf("Error unmarshalling response: %v", err)
	}
	w.broadcasterID = entity.Data[0].Id
	viper.Set("twitchID", w.broadcasterID)
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config: %v", err)
	}

	return nil
}

func (w *WebSocket) Connect() {
	var err error
	w.conn, _, err = websocket.DefaultDialer.Dial(w.url, nil)
	if err != nil {
		log.Fatalf("error connecting to %s: %s", w.url, err)
	}
	defer func(Conn *websocket.Conn) {
		err := Conn.Close()
		if err != nil {
			log.Error("close:", err)
		}
	}(w.conn)

	for {
		_, message, err := w.conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			continue
		}

		var twitchMsg TwitchMessage
		if err := json.Unmarshal(message, &twitchMsg); err != nil {
			log.Error("unmarshal twitchMsg:", err)
			continue
		}
		switch twitchMsg.Metadata.MessageType {
		case "session_welcome":
			var sessionWelcomeMessage SessionWelcomeMessage
			if err := json.Unmarshal(message, &sessionWelcomeMessage); err != nil {
				log.Error("unmarshal sessionWelcomeMessage:", err)
				continue
			}
			sessionID := sessionWelcomeMessage.Payload.Session.ID
			log.Debug("sessionID: " + sessionID)
			w.createSubscription(sessionID)
		case "notification":
			var rewardRedeemedMessage RewardRedeemedMessage
			if err := json.Unmarshal(message, &rewardRedeemedMessage); err != nil {
				log.Error("unmarshal chatMessage:", err)
				continue
			}
			log.Debug("Сработал скип трека: " + rewardRedeemedMessage.Payload.Event.Reward.Title)
			if rewardRedeemedMessage.Payload.Event.Reward.Title == "Скип трека" {
				w.Repository.Skip()
			}
		}
	}
}

func (w *WebSocket) createSubscription(sessionID string) {
	broadcasterUserID := viper.GetString("twitchID")

	requestBody := map[string]any{
		"type":    "channel.channel_points_custom_reward_redemption.add",
		"version": "1",
		"condition": map[string]string{
			"broadcaster_user_id": broadcasterUserID,
		},
		"transport": map[string]string{
			"method":     "websocket",
			"session_id": sessionID,
		},
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatalf("Error marshaling request body: %v", err)
	}

	req, resp := httpClient.GetRequest(w.tokenModel.AccessToken)
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://api.twitch.tv/helix/eventsub/subscriptions")
	req.Header.SetMethod("POST")
	req.SetBody(jsonBody)
	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	req.Header.Set("Client-Id", w.clientID)
	req.Header.Set("Content-Type", "application/json")

	err = w.httpClient.Do(req, resp)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	log.Debug(req)

	if resp.StatusCode() != fasthttp.StatusAccepted {
		log.Debug(resp)
		log.Fatalf("Error: received non-202 response code: %v", resp.StatusCode())
	}

	log.Debug("Subscription created successfully")
}
