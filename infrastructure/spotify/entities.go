package spotify

import (
	"fmt"
	"test/test/domain/models"
)

type userResponse struct {
	UserName string `json:"display_name"`
}

type searchResponse struct {
	Tracks struct {
		Items []struct {
			Name    string `json:"name"`
			URI     string `json:"uri"`
			Artists []struct {
				Name string `json:"name"`
			} `json:"artists"`
		} `json:"items"`
	} `json:"tracks"`
}

func (r *searchResponse) MakeModel() models.SearchTrack {
	return models.SearchTrack{
		TrackName: r.Tracks.Items[0].Name,
		TrackURI:  r.Tracks.Items[0].URI,
		Artist:    r.Tracks.Items[0].Artists[0].Name,
	}
}

func (r *userResponse) MakeModel() models.User {
	return models.User{Username: r.UserName}
}

type trackResponse struct {
	Item struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"item"`
}

type queueResponse struct {
	CurrentlyPlaying struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"currently_playing"`
	Queue []struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"queue"`
}

type errorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (q *queueResponse) FeaturesTracks() string {
	var queue string
	var artists string
	for i := range q.Queue[0].Artists {
		artists += q.Queue[0].Artists[i].Name
		if i != len(q.Queue[0].Artists)-1 {
			artists += ", "
		}
	}
	for i := 0; i < 5; i++ {
		queue = fmt.Sprint(queue, fmt.Sprintf("%d. %s - %s, ", i+1, q.Queue[i].Artists[0].Name, q.Queue[i].Name))
	}
	queue = fmt.Sprintf("В очереди: %s.", queue[:len(queue)-2])
	return queue
}
