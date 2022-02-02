package util

import (
	"context"
	"fmt"
	"github.com/ably-labs/word-game/word-game-be/model"
	"github.com/ably/ably-go/ably"
)

func PublishLobbyMessage(client *ably.Realtime, lobby *int64, name string, message interface{}) error {
	return client.Channels.Get(fmt.Sprintf("lobby-%d", *lobby)).Publish(context.Background(), name, message)
}

// LobbyListUpdate updates the lobby list for public lobbies when a lobby state chanes
func LobbyListUpdate(client *ably.Realtime, lobby *model.Lobby) error {
	if lobby.Private {
		return nil

	}
	return client.Channels.Get("lobby-list").Publish(context.Background(), "lobbyUpdate", lobby)
}
