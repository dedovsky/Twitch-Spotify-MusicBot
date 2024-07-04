package spotify

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

type Repository struct {
	clientID      string
	Base64Secrets string

	client       *fasthttp.Client
	tokenModel   models.Token
	TokenChannel chan string
	redirectURI  string
}

func NewRepository() *Repository {
	r := &Repository{
		clientID:      "",
		Base64Secrets: "", //кодированный Base64 client_id:client_secret

		client:       &fasthttp.Client{},
		TokenChannel: make(chan string),
		redirectURI:  "http://localhost:8080/callback",
		tokenModel: models.Token{
			RefreshToken: viper.GetString("refreshToken"),
		},
	}
	r.SetupToken()
	return r
}

func (r *Repository) SetupToken() {
	if r.tokenModel.RefreshToken != "" {
		log.Debug("Токен найден")
		r.refreshToken()
		log.Info("Токен обновлен")
		return
	}
	ser := http.NewService()
	go ser.StartServer()
	clientId := r.clientID
	log.Info("Ссылка для логина: " + ("https://accounts.spotify.com/authorize?response_type=code&client_id=" +
		clientId +
		"&scope=user-read-playback-state%20user-read-currently-playing%20user-modify-playback-state&redirect_uri=" +
		r.redirectURI))
	log.Debug("rToken: " + r.tokenModel.RefreshToken)

	code := <-ser.GetChan()
	ser.CloseChan()
	ser.StopServer()

	log.Debug("Код получен")
	r.GetTokenFromCode(code)
	log.Debug("Токен установлен")
	go tokenService.RefreshEveryExpiresIn(r)
}

func (r *Repository) refreshToken() {
	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://accounts.spotify.com/api/token")
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Set("Authorization", "Basic "+r.Base64Secrets)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {r.tokenModel.RefreshToken},
	}
	req.SetBodyString(data.Encode())
	if err := r.client.Do(req, resp); err != nil {
		log.Error("ошибка получения токена: ", err)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		// если описание ошибки Refresh Base64Secrets revoked
		var errorResp errorResponse
		_ = json.Unmarshal(resp.Body(), &errorResp)
		if errorResp.ErrorDescription == "Refresh Base64Secrets revoked" {
			viper.Set("refreshToken", "")
			r.SetupToken()
		}
		log.Debug("Body: " + string(resp.Body()))
		log.Error("Ошибка при получении токена спотифая: " + errorResp.ErrorDescription)
	}
	var entity models.Token
	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		log.Error("Ошибка десериализации: ", err)
	}
	r.tokenModel = entity
	if entity.RefreshToken != "" {
		viper.Set("refreshToken", entity.RefreshToken)
		err := viper.WriteConfig()
		if err != nil {
			log.Error("Ошибка записи в конфиг: ", err)
		}
	}
}

func (r *Repository) GetTokenFromCode(code string) {
	req, resp := httpClient.GetRequest("")
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://accounts.spotify.com/api/token")
	req.Header.SetMethod(fasthttp.MethodPost)

	req.Header.Set("Authorization", "Basic "+r.Base64Secrets)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {"http://localhost:8080/callback"},
	}
	req.SetBodyString(data.Encode())
	if err := r.client.Do(req, resp); err != nil {
		log.Error("Error: ", err)
	}
	var entity models.Token
	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		log.Error("Ошибка десериализации: ", err)
	}
	r.tokenModel = entity
	log.Debug("Токен: " + entity.AccessToken)
	log.Debug("Рефреш токен: " + entity.RefreshToken)
	viper.Set("refreshToken", entity.RefreshToken)
	err := viper.WriteConfig()
	if err != nil {
		log.Error("Ошибка записи в конфиг: ", err)
	}
}

func (r *Repository) GetExpiresIn() int {
	return r.tokenModel.ExpiresIn
}
