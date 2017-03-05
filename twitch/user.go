package twitch

type user struct {
	Username string
	UserId string
	Color string
	DisplayName string
	Emotes map[string][]string
	Mod bool
	Subscriber bool
	Turbo bool
}

func newUser(username string) user {
	return user{
		Username: username,
	}
}