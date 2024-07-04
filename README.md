# Twitch-Spotify-MusicBot

Twitch музыкальный бот, подключающийся к Spotify. 
Был создан по просьбе человека с хардкоднутыми токенами и прочим. По сооброжениям безопасности, их пришлось убрать. Чтобы добавить свои - нужно отредактировать структуры по путям:

```
services/twitch/service.go
services/websocket/websocket.go
infrastructure/spotify/repository.go
```

Доступные команды в боте:
```
!ping
!song
!sr <название/ссылка на трек>
!queue
```
