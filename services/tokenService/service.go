package tokenService

import (
	"time"
)

type serviceInterface interface {
	GetExpiresIn() int
	SetupToken()
}

func RefreshEveryExpiresIn(s serviceInterface) {
	for {
		time.Sleep(time.Duration(s.GetExpiresIn())*time.Second - 1*time.Minute)
		s.SetupToken()
	}
}
