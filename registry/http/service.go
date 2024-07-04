package http

import (
	"github.com/charmbracelet/log"
	"github.com/valyala/fasthttp"
)

type Service struct {
	server   *fasthttp.Server
	url      string
	token    string
	code     string
	codeChan chan string
}

func NewService() *Service {
	service := &Service{
		server:   &fasthttp.Server{},
		codeChan: make(chan string),
	}
	service.server.Handler = func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/callback":
			log.Debug("Сработал callback")
			service.handleCallback(ctx)
		default:
			log.Debug("Сработал default")
			service.handleRequest(ctx)

		}
	}
	return service
}

func (s *Service) StartServer() {
	log.Debug("Сервер запущен")
	if err := s.server.ListenAndServe(":8080"); err != nil {
		log.Error("Error: " + err.Error())
	}
	log.Debug("Сервер выключен")

}

func (s *Service) StopServer() {
	err := s.server.Shutdown()
	if err != nil {
		log.Error("Ошибка выключения: " + err.Error())
	}
}

func (s *Service) GetChan() chan string {
	return s.codeChan
}

func (s *Service) CloseChan() {
	close(s.codeChan)
}
