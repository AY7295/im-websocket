# websocket

## GET /chat?receiver=[id]
注: 不在线用户的消息 会由极光推送
### header
```
{
  "Authorization": "<token>"
}
```
### message json format
```json5
{
    "_id": "1",
    "text": "这里是要发送信息",
    "createdAt": "2022-07-24T18:19:44.645+08:00",
    "user": {
        "id": "发送者的id",
        "name": "发送者的姓名",
        "avatar": "发送者 头像 url"
    },
    "image": "img url",
    "sent": false,
    "received": false,
    "pending": false
}
```