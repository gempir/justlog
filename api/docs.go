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

//swagger:parameters list
type ListParams struct {
	// in: query
	// required: true
	Channelid string `json:"channelid"`
	// in: query
	// required: true
	Userid string `json:"userid"`
	LogParams
}

//swagger:parameters userLogs
type UserLogsParams struct {
	// in: path
	Channel string `json:"channel"`
	// in: path
	Username string `json:"username"`
	LogParams
}

//swagger:parameters userLogsYearMonth
type UserLogsYearMonthParams struct {
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
