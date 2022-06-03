package model

type RequestBody struct {
	ParentUuid string `json:"parent_uuid"`
	Query      string `json:"query"`
	Type       string `json:"type"`
}

type ResponseBody struct {
	Items []struct {
		Additional string `json:"additional"`
		Caption    string `json:"caption"`
		Type       string `json:"type"`
		Uuid       string `json:"uuid"`
	} `json:"items"`
	Total int `json:"total"`
}

type BitopBody struct {
	Token   string
	Request RequestBody
}
