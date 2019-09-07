package filelog

import (
	"strings"
	"time"
)

func logChannelIdentifier(channel string, roomID string) string {
	if strings.HasPrefix(channel, "chatrooms:") {
		return channel
	} else {
		return roomID
	}
}

func getLogTime(msgChannel string, msgTime *time.Time) (int, int, int) {
	// this exists because the `tmi-sent-ts` tag on messages from chatrooms can be invalid
	// (at the time of this patch, which is 2019-09-07)
	// We would receive a negative value that represents some time in the 18th century
	// TODO: Once twitch has fixed this behaviour, remove this conditional, and just do
	//       what the `else` block does
	if strings.HasPrefix(msgChannel, "chatrooms:") {
		year, month, day := time.Now().Date()
		return year, int(month), day
	} else {
		year, month, day := msgTime.Date()
		return year, int(month), day
	}
}
