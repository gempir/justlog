package twitch

type user struct {
	Username string
	UserId string
	Color string
	DisplayName string
	Mod bool
	Subscriber bool
	Turbo bool
	Emotes map[string][]string
}

func newUser(username string, userId string, color string, displayName string, mod bool, subscriber bool, turbo bool) user {
	return user{
		Username: username,
		UserId: userId,
		Color: color,
		DisplayName: displayName,
		Mod: mod,
		Subscriber: subscriber,
		Turbo: turbo,
	}
}