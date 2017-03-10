package twitch

import (
	"strings"
	"strconv"
)

func parseMessage(msg string) Message {

	fullUser := userReg.FindString(msg)
	userIrc := strings.Split(fullUser, "!")
	username := userIrc[0][1:len(userIrc[0])]
	split2 := strings.Split(msg, ".tmi.twitch.tv PRIVMSG ")
	split3 := channelReg.FindString(split2[1])
	channel := split3[0 : len(split3)-2]
	split4 := strings.Split(split2[1], split3)
	message := split4[1]
	message = actionReg.ReplaceAllLiteralString(message, "")
	message = actionReg2.ReplaceAllLiteralString(message, "")

	split5 := strings.Replace(split2[0], " :" + username + "!" + username + "@" + username, "", -1)
	tags := strings.Split(strings.Replace(split5, "@", "", 1), ";")

	tagMap := make(map[string]string)

	for _,tag := range tags {
		tagSplit := strings.Split(tag, "=")
		tagMap[tagSplit[0]] = tagSplit[1]
	}

	subscriber, turbo, mod := false, false, false

	if tagMap["subscriber"] == "1" {
		subscriber = true
	}
	if tagMap["turbo"] == "1" {
		turbo = true
	}
	if tagMap["mod"] == "1" {
		mod = true
	}

	emotes := getEmotes(message, tagMap["emotes"])


	user := newUser(username, tagMap["User-id"], tagMap["color"], tagMap["display-name"], mod, subscriber, turbo, emotes)
	return newMessage(message, user, channel)
}


func getEmotes(message string, emotesString string) []Emote {
	emotes := make([]Emote, 0)

	if emotesString == "" {
		return emotes
	}

	emoteStrSplit := strings.Split(emotesString, "/")

	for _, emoteStr := range emoteStrSplit {

		emoteSplit := strings.Split(emoteStr, ":")

		emotePositions := strings.Split(emoteSplit[1], ",")
		emoteCount := len(emotePositions)

		emotePosition := strings.Split(emotePositions[0], "-")

		pos1, err1 := strconv.Atoi(emotePosition[0])
		pos2, err2 := strconv.Atoi(emotePosition[1])

		if err1 != nil || err2 != nil {
			continue
		}

		emoteCode := message[pos1:pos2+1]

		for i := 0; i < emoteCount; i++ {

			id, err := strconv.Atoi(emoteSplit[0])
			if err != nil {
				continue
			}
			emote := NewEmote(id, emoteCode)
			emotes = append(emotes, emote)
		}
	}

	return emotes
}