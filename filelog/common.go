package filelog

import "strings"

func logChannelIdentifier(channel string, roomID string) string {
	if strings.HasPrefix(channel, "chatrooms:") {
		return channel
	} else {
		return roomID
	}
}

