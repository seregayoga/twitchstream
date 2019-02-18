package twitch

import (
	"encoding/json"
	"fmt"

	twitchirc "github.com/gempir/go-twitch-irc"
	"github.com/gorilla/websocket"
)

// Event stream event
type Event struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

// IRCToWSStreamer streamer for events from irc to websocket
type IRCToWSStreamer struct {
	ws           *websocket.Conn
	ircClient    *twitchirc.Client
	streamerName string
	token        string
}

// NewIRCToWSStreamer creates new IRCToWSStreamer
func NewIRCToWSStreamer(ws *websocket.Conn, userName string, streamerName string, token string) *IRCToWSStreamer {
	return &IRCToWSStreamer{
		ws:           ws,
		ircClient:    twitchirc.NewClient(userName, "oauth:"+token),
		streamerName: streamerName,
		token:        token,
	}
}

// Stream streams events from irc to websocket
func (ircws *IRCToWSStreamer) Stream() error {
	ircws.ircClient.OnNewWhisper(func(user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "whisper",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewRoomstateMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "room state message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewClearchatMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "clear chat message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewUsernoticeMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "user notice message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewNoticeMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "notice message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnNewUserstateMessage(func(channel string, user twitchirc.User, message twitchirc.Message) {
		event := Event{
			Type:    "user state message",
			Content: fmt.Sprintf("%s: %s", user.DisplayName, message.Text),
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnUserJoin(func(channel, user string) {
		event := Event{
			Type:    "user join",
			Content: user,
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.OnUserPart(func(channel, user string) {
		event := Event{
			Type:    "user part",
			Content: user,
		}

		data, _ := json.Marshal(event)

		err := ircws.ws.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			ircws.ircClient.Disconnect()
		}
	})

	ircws.ircClient.Join(ircws.streamerName)

	return ircws.ircClient.Connect()
}
