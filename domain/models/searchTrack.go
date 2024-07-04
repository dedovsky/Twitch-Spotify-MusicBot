package models

type SearchTrack struct {
	TrackName string `json:"trackName"`
	TrackURI  string `json:"trackURI"`
	Artist    string `json:"artist"`
}
