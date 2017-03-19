package twitch

type Command struct {
	IsCommand bool
	Name string
	Args []string
}