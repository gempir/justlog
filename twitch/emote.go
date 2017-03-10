package twitch

type Emote struct {
	Id   int
	Code string
}

func NewEmote(id int, code string) Emote {
	return Emote{
		Id:   id,
		Code: code,
	}
}