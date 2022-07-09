# websocket

## GET /ping
### response
```
{
  "message": "pong"
}
```

## GET /chat
注: 不在线用户的消息 会由极光推送
### header
```
{
  "Authorization": "<token>"
}
```