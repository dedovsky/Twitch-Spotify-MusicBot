package websocket

import (
	"encoding/json"
	"github.com/charmbracelet/log"
	"github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"net/url"
	"test/test/domain/models"
	"test/test/infrastructure/httpClient"
	"test/test/registry/http"
	"test/test/services/tokenService"
)

func (w *WebSocket) SetupToken() {
	w.tokenModel.RefreshToken = viper.GetString("ChannelToken")
	if w.tokenModel.RefreshToken != "" {
		log.Debug("Токен найден")
		w.refreshToken()
		log.Info("Токен обновлен")
		return
	}
	ser := http.NewService()
	go ser.StartServer()
	link := "https://id.twitch.tv/oauth2/authorize?response_type=code&client_id=" +
		w.clientID +
		"&redirect_uri=http://localhost:8080/callback&scope=chat:read+chat:edit+user:read:email+channel:read:redemptions"
	log.Info("Ссылка для логина: " + link)

	code := <-ser.GetChan()
	ser.CloseChan()
	ser.StopServer()

	log.Debug("Код получен")
	w.GetTokenFromCode(code)
	log.Debug("Токен установлен: " + w.tokenModel.RefreshToken)
	go tokenService.RefreshEveryExpiresIn(w)

}

func (w *WebSocket) refreshToken() {
	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://id.twitch.tv/oauth2/token")
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data := url.Values{
		"client_id":     {w.clientID},
		"client_secret": {w.clientSecret},
		"grant_type":    {"refresh_token"},
		"refresh_token": {w.tokenModel.RefreshToken},
	}
	req.SetBodyString(data.Encode())
	if err := w.httpClient.Do(req, resp); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		viper.Set("ChannelToken", "")
		log.Debug("Body: " + string(resp.Body()))
		log.Fatalf("Ошибка при обновлении токена вебхука: %v", resp.StatusCode())
	}
	var entity models.Token
	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		log.Fatalf("Error: %v", err)
	}
	w.tokenModel = entity
	if w.tokenModel.RefreshToken != "" {
		viper.Set("ChannelToken", w.tokenModel.RefreshToken)
		err := viper.WriteConfig()
		if err != nil {
			log.Errorf("Error writing config: %v", err)
		}
	}

}

func (w *WebSocket) GetTokenFromCode(code string) {
	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://id.twitch.tv/oauth2/token")
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data := url.Values{
		"client_id":     {w.clientID},
		"client_secret": {w.clientSecret},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {"http://localhost:8080/callback"},
		"code":          {code},
	}
	req.SetBodyString(data.Encode())
	if err := w.httpClient.Do(req, resp); err != nil {
		log.Fatalf("Error: %v", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		log.Debug("Body: " + string(resp.Body()))
		log.Fatalf("Error: received non-200 response code: %s", resp.StatusCode())
	}
	var entity models.Token
	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		log.Fatalf("Error: %v", err)
	}
	w.tokenModel = entity
	viper.Set("ChannelToken", w.tokenModel.RefreshToken)
	err := viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config: %v", err)
	}
}

func (w *WebSocket) GetExpiresIn() int {
	return w.tokenModel.ExpiresIn
}
