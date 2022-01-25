package util

import (
	"context"
	"fmt"
	"github.com/ably/ably-go/ably"
)

func PublishLobbyMessage(client *ably.Realtime, lobby *int64, name string, message interface{}) error {
	return client.Channels.Get(fmt.Sprintf("lobby-%d", *lobby)).Publish(context.Background(), name, message)
}
