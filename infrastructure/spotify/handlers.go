package spotify

import (
	"encoding/json"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/valyala/fasthttp"
	"net/url"
	"test/test/domain/models"
	"test/test/infrastructure/httpClient"
	"time"
)

func (r *Repository) Ping() time.Duration {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	defer httpClient.ReleaseRR(req, resp)

	req.SetRequestURI("https://api.spotify.com/v1/me")
	start := time.Now()

	if err := r.client.Do(req, resp); err != nil {
		log.Error("Error: ", err)
	}
	end := time.Since(start)
	return end
}

func (r *Repository) Song() string {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	req.SetRequestURI("https://api.spotify.com/v1/me/player/currently-playing")
	defer httpClient.ReleaseRR(req, resp)

	if err := r.client.Do(req, resp); err != nil {
		log.Error("Error: ", err)
	}

	var aTrack trackResponse
	if err := json.Unmarshal(resp.Body(), &aTrack); err != nil {
		log.Error("Ошибка десериализации: ", err)
		log.Debug(resp)
		return "Ошибка подключения Spotify"
	}
	if len(aTrack.Item.Artists) < 0 {
		return "Сейчас ничего не играет"
	}
	var artists string
	for i := range aTrack.Item.Artists {
		artists += aTrack.Item.Artists[i].Name
		if i != len(aTrack.Item.Artists)-1 {
			artists += ", "
		}
	}
	track := fmt.Sprintf("Сейчас играет: %s - %s", artists, aTrack.Item.Name)

	return track

}

func (r *Repository) Sr(song string) (models.SearchTrack, error) {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	defer httpClient.ReleaseRR(req, resp)
	req.SetRequestURI("https://api.spotify.com/v1/search?q=" + url.QueryEscape(song) + "&type=track")

	if err := r.client.Do(req, resp); err != nil {
		log.Error("Ошибка запроса: ", err)
	}
	if resp.StatusCode() != 200 {
		log.Error("Ошибка поиска: ", resp.BodyStream())
		return models.SearchTrack{}, fmt.Errorf("ошибка")
	}

	var search searchResponse
	if err := json.Unmarshal(resp.Body(), &search); err != nil {
		log.Error("Ошибка десериализации во время поиска: ", err)
	}
	if len(search.Tracks.Items) == 0 {
		return models.SearchTrack{}, fmt.Errorf("трек не найден")
	}
	err := r.srPlay(search.Tracks.Items[0].URI)
	if err != nil {
		return models.SearchTrack{}, fmt.Errorf("ошибка воспроизведения: %w", err)
	}
	return search.MakeModel(), nil
}

func (r *Repository) srPlay(song string) error {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	defer httpClient.ReleaseRR(req, resp)
	req.SetRequestURI("https://api.spotify.com/v1/me/player/queue?uri=" + url.QueryEscape(song))
	req.Header.SetMethod(fasthttp.MethodPost)
	if err := r.client.Do(req, resp); err != nil {
		return err
	}
	if resp.StatusCode() != 204 {
		log.Error("Error: ", resp.BodyStream())
		return fmt.Errorf("неизвестная ошибка")
	}
	return nil
}

func (r *Repository) GetUser() models.User {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	defer httpClient.ReleaseRR(req, resp)
	req.SetRequestURI("https://api.spotify.com/v1/me")

	if err := r.client.Do(req, resp); err != nil {
		return models.User{}
	}

	var entity userResponse

	if err := json.Unmarshal(resp.Body(), &entity); err != nil {
		return models.User{}
	}

	return entity.MakeModel()
}

func (r *Repository) GetQueue() string {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	req.SetRequestURI("https://api.spotify.com/v1/me/player/queue")
	defer httpClient.ReleaseRR(req, resp)

	if err := r.client.Do(req, resp); err != nil {
		log.Error("Error: ", err)
	}

	var aQueue queueResponse
	if err := json.Unmarshal(resp.Body(), &aQueue); err != nil {
		log.Error("Error: ", err)
	}
	if len(aQueue.Queue) == 0 {
		return "Сейчас ничего нет в очереди"
	}
	var artists string
	for i := range aQueue.CurrentlyPlaying.Artists {
		artists += aQueue.CurrentlyPlaying.Artists[i].Name
		if i != len(aQueue.CurrentlyPlaying.Artists)-1 {
			artists += ", "
		}
	}

	queue := fmt.Sprintf(
		"Сейчас играет: %s - %s. %s",
		artists,
		aQueue.CurrentlyPlaying.Name,
		aQueue.FeaturesTracks(),
	)
	return queue
}

func (r *Repository) Skip() {
	req, resp := httpClient.GetRequest(r.tokenModel.AccessToken)
	req.SetRequestURI("https://api.spotify.com/v1/me/player/next")
	req.Header.SetMethod(fasthttp.MethodPost)
	defer httpClient.ReleaseRR(req, resp)

	req.Header.Set("Authorization", "Bearer "+r.tokenModel.AccessToken)

	if err := r.client.Do(req, resp); err != nil {
		log.Error("Error: ", err)
	}

}
