package http

import (
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/valyala/fasthttp"
)

func (s *Service) handleCallback(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(200)
	ctx.SetContentType("text/html; charset=utf-8")

	code := string(ctx.QueryArgs().Peek("code"))
	_, _ = fmt.Fprintf(ctx, "Токен получен. Откройте приложение.")
	s.codeChan <- code
	log.Debug("Code: " + code)
	return
}

func (s *Service) handleRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(404)
	ctx.SetContentType("text/html; charset=utf-8")
	_, _ = fmt.Fprintf(ctx, "Скорее всего произошла ошибка.")

}
