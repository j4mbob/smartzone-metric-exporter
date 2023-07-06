package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"regexp"
)

type Alarms struct {
	List []struct {
		ID            string `json:"id"`
		Acknowledged  string `json:"acknowledged"`
		Activity      string `json:"activity"`
		AlarmCode     int64  `json:"alarmCode"`
		AlarmState    string `json:"alarmState"`
		AlarmType     string `json:"alarmType"`
		Category      string `json:"category"`
		InsertionTime int64  `json:"insertionTime"`
		Severity      string `json:"severity"`
	} `json:"list"`
	TotalCount int64 `json:"totalCount"`
}

func GetActiveAlarms(controller map[string]string, domainsData []Domains) map[string]int {

	radiusAlarmQuery := BuildAlarmQuery("Radius server unreachable outstanding", []string{"alarmType", "alarmState"})

	activeAlarms := GetAlarmsQuery(controller, radiusAlarmQuery)

	return ParseActiveAlarms(activeAlarms, domainsData)

	//return GetActiveAlarmsTotal((activeAlarms))

}

func BuildAlarmQuery(searchText string, searchFields []string) []byte {

	type Filter struct {
		Type     string `json:"type"`
		Value    string `json:"value"`
		Operator string `json:"operator"`
	}

	type TextSearch struct {
		Fields []string `json:"fields"`
		Type   string   `json:"type"`
		Value  string   `json:"value"`
	}

	type QueryFilter struct {
		FilterList []Filter   `json:"filters"`
		TextSearch TextSearch `json:"fullTextSearch"`
		Limit      int64      `json:"limit"`
	}

	filter := QueryFilter{
		FilterList: []Filter{{Type: "CATEGORY", Value: "AP", Operator: "eq"}}, Limit: 1000}

	filter.TextSearch.Type = "AND"

	filter.TextSearch.Value = searchText
	filter.TextSearch.Fields = searchFields

	queryFilter, err := json.Marshal(filter)

	if err != nil {
		fmt.Println("error constructing query filter: ", err)
	}

	return queryFilter

}

func GetAlarmsQuery(controller map[string]string, queryFilter []byte) Alarms {

	queryApUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/alert/alarm/list"

	httpResp, _ := BuildHttpRequest(queryApUrl, "POST", nil, bytes.NewBuffer(queryFilter), controller["accesstoken"], true)

	var alarms Alarms

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &alarms)
	if err != nil {
		fmt.Println(err)
	}

	return alarms

}

func GetActiveAlarmsTotal(alarms Alarms) int {

	return len(alarms.List)
}

func ParseActiveAlarms(alarms Alarms, domainsData []Domains) map[string]int {

	radiusServers := make(map[string]int)

	for _, v := range alarms.List {
		var re = regexp.MustCompile(`(?m)\[RuckusAP@(?P<mac>.*?)\].*\[(?P<ip>.*?)\]`)
		match := re.FindStringSubmatch(v.Activity)
		//apMac := re.SubexpIndex("mac")
		ipIndex := re.SubexpIndex("ip")
		radiusIP := match[ipIndex]

		_, siteName := GetSiteName(radiusIP, domainsData)

		radiusServers[siteName] += 1

	}
	return radiusServers

}
