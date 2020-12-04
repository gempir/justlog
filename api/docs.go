package api

type LogParams struct {
	// in: query
	Json string `json:"json"`
	// in: query
	Reverse string `json:"reverse"`
	// in: query
	From int32 `json:"from"`
	// in: query
	To int32 `json:"to"`
}

//swagger:parameters channelUserLogsRandom
type ChannelUserLogsRandomParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Username string `json:"username"`
	LogParams
}

//swagger:parameters channelUserLogs
type ChannelUserLogsParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Username string `json:"username"`
	LogParams
}

//swagger:parameters channelUserLogsYearMonth
type ChannelUserLogsYearMonthParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Username string `json:"username"`
	// in: path
	Year string `json:"year"`
	// in: path
	Month string `json:"month"`
	LogParams
}

//swagger:parameters channelLogs
type ChannelLogsParams struct {
	// in: path
	Channel string `json:"channel"`
	LogParams
}

//swagger:parameters channelLogsYearMonthDay
type ChannelLogsYearMonthDayParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Year string `json:"year"`
	// in: path
	Month string `json:"month"`
	// in: path
	Day string `json:"day"`
	LogParams
}

//swagger:parameters channelIdUserIdLogsRandom
type ChannelIdUserIdLogsRandomParams struct {
	// in: path
	ChannelId string `json:"channelid"`
	// in: path
	UserId string `json:"userid"`
	LogParams
}

//swagger:parameters channelIdUserIdLogs
type ChannelIdUserIdLogsParams struct {
	// in: path
	ChannelId string `json:"channelid"`
	// in: path
	UserId string `json:"userid"`
	LogParams
}

//swagger:parameters channelIdUserIdLogsYearMonth
type ChannelIdUserIdLogsYearMonthParams struct {
	// in: path
	ChannelId string `json:"channelid"`
	// in: path
	UserId string `json:"userid"`
	// in: path
	Year string `json:"year"`
	// in: path
	Month string `json:"month"`
	LogParams
}

//swagger:parameters channelIdLogs
type ChannelIdLogsParams struct {
	// in: path
	Channel string `json:"channelid"`
	LogParams
}

//swagger:parameters channelIdLogsYearMonthDay
type ChannelIdLogsYearMonthDayParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Year string `json:"year"`
	// in: path
	Month string `json:"month"`
	// in: path
	Day string `json:"day"`
	LogParams
}
