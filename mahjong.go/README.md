# Parlour
Mahjong Party backend

## API

### Host game

* Method: `POST`
* Path: `/rooms`
* Headers:
  * Content-Type: `application/x-www-form-urlencoded`
* Body: `name=:name`

Returns the ID of the newly-created room.

### Join game

* Method: `POST`
* Path: `/rooms/:id/players`
* Headers:
  * Content-Type: `application/x-www-form-urlencoded`
* Body: `name=:name`

### Subscribe to game updates

Path: `/rooms/:id/live`

This in an [SSE](https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events) endpoint. Use an [EventSource](https://developer.mozilla.org/en-US/docs/Web/API/EventSource) to consume it.

Each message will be a JSON-encoded `RoomView` struct.

### Do something (draw, discard, chi, pong etc.)

* Method: `POST`
* Path: `/rooms/:id/actions`
* Headers:
  * Content-Type: `application/json`
* Body: `{"type": "chi", "data": {"tiles": ["22一索","23二索"]}}`

If the action is successful, the updated game state will be broadcast to connected clients.
