package api

//swagger:parameters list
type ListParams struct {
	// in: query
	// required: true
	Channelid string `json:"channelid"`
	// in: query
	// required: true
	Userid string `json:"userid"`
}
