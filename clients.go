package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func GetClientCountsDomain(controller map[string]string, domainData []Domains) map[string]int {

	domainClientCounts := make(map[string]int)

	for _, k := range domainData {

		clientCount := GetClientCount(controller, k.ID)
		domainClientCounts[k.Name] = clientCount
	}

	return domainClientCounts

}

func GetClientCount(controller map[string]string, searchDomain string) int {

	type ClientsQuery struct {
		TotalCount int `json:"totalCount"`
	}

	var queryFilter []byte

	if searchDomain == "" {
		queryFilter = BuildGetFullClientCountQuery()

	} else {
		queryFilter = BuildClientCountsDomainQuery(searchDomain)
	}

	queryClientsTotalUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/query/client"
	httpResp, _ := BuildHttpRequest(queryClientsTotalUrl, "POST", nil, bytes.NewBuffer(queryFilter), controller["accesstoken"], true)

	var clientsQuery ClientsQuery

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &clientsQuery)
	if err != nil {
		fmt.Println(err)
	}

	return clientsQuery.TotalCount

}

func BuildClientCountsDomainQuery(searchDomain string) []byte {

	type Filter struct {
		Type     string `json:"type"`
		Value    string `json:"value"`
		Operator string `json:"operator"`
	}

	type TextSearch struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type QueryFilter struct {
		FilterList []Filter   `json:"filters"`
		TextSearch TextSearch `json:"fullTextSearch"`
		Limit      int64      `json:"limit"`
	}
	textSearch := TextSearch{Type: "AND", Value: ""}

	filter := QueryFilter{
		FilterList: []Filter{{Type: "DOMAIN", Value: searchDomain, Operator: "eq"}}, TextSearch: textSearch, Limit: 1000}

	queryFilter, err := json.Marshal(filter)

	if err != nil {
		fmt.Println("error constructing query filter: ", err)
	}

	return queryFilter

}

func BuildGetFullClientCountQuery() []byte {

	type TextSearch struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type QueryFilter struct {
		TextSearch TextSearch `json:"fullTextSearch"`
		Limit      int64      `json:"limit"`
	}

	textSearch := TextSearch{Type: "AND", Value: ""}
	filter := QueryFilter{Limit: 1000, TextSearch: textSearch}

	queryFilter, err := json.Marshal(filter)

	if err != nil {
		fmt.Println("error constructing query filter: ", err)
	}

	return queryFilter

}
