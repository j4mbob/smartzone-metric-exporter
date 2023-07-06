package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func GetApTotals(controller map[string]string) (int, int) {

	onlineQueryFilter := BuildApQuery("online")
	offlineQueryFilter := BuildApQuery("offline")

	return GetApQuery(controller, onlineQueryFilter), GetApQuery(controller, offlineQueryFilter)

}

func BuildApQuery(queryType string) []byte {

	type Filter struct {
		Type     string `json:"type"`
		Value    string `json:"value"`
		Operator string `json:"operator"`
	}

	type QueryFilter struct {
		FilterList []Filter `json:"filters"`
	}

	var filter QueryFilter

	if queryType == "online" {
		filter = QueryFilter{
			FilterList: []Filter{{Type: "SYNCEDSTATUS", Value: "Online", Operator: "eq"}}}
	} else if queryType == "offline" {
		filter = QueryFilter{
			FilterList: []Filter{{Type: "SYNCEDSTATUS", Value: "Offline", Operator: "eq"}}}
	}

	queryFilter, err := json.Marshal(filter)

	if err != nil {
		fmt.Println("error constructing query filter: ", err)
	}

	return queryFilter

}

func GetApQuery(controller map[string]string, queryFilter []byte) int {

	type ApCount struct {
		TotalCount int `json:"totalCount"`
	}
	queryApUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/query/ap"

	httpResp, _ := BuildHttpRequest(queryApUrl, "POST", nil, bytes.NewBuffer(queryFilter), controller["accesstoken"], true)

	var ApTotal ApCount

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &ApTotal)
	if err != nil {
		fmt.Println("error: ", err)
	}

	return ApTotal.TotalCount
}
