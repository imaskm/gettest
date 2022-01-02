package types

import "time"

type Record struct {
	Key        string    `json:"key" bson:"key"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	TotalCount int       `json:"totalCount" bson:"totalCount"`
}

type RecordDb struct {
	Key       string    `json:"key" bson:"key"`
	Value     string    `json:"value" bson:"value"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Counts    []int     `json:"totalCount" bson:"counts"`
}

type RecordRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	MinCount  int    `json:"minCount"`
	MaxCount  int    `json:"maxCount"`
}

type RecordResponse struct {
	Code    int       `json:"code"`
	Msg     string    `json:"msg"`
	Records []*Record `json:"records"`
}

type Data struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
