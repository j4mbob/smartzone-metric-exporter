package main

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func GetZoneRadiusServer(controller map[string]string, zoneId string) []Radius {

	type RadiusServersQuery struct {
		List []struct {
			ID          string `json:"id"`
			Description string `json:"description"`
			Name        string `json:"name"`
			Primary     struct {
				IP           string `json:"ip"`
				Port         int64  `json:"port"`
				SharedSecret string `json:"sharedSecret"`
			} `json:"primary"`
			Secondary struct {
				IP           string `json:"ip"`
				Port         int64  `json:"port"`
				SharedSecret string `json:"sharedSecret"`
			} `json:"secondary"`
			ServiceType string `json:"serviceType"`
			ZoneID      string `json:"zoneId"`
		} `json:"list"`
	}

	queryZoneRadiusServerUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/rkszones/" + zoneId + "/aaa/radius"

	httpResp, _ := BuildHttpRequest(queryZoneRadiusServerUrl, "GET", nil, nil, controller["accesstoken"], true)

	var radiusServersQuery RadiusServersQuery

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &radiusServersQuery)
	if err != nil {
		fmt.Println(err)
	}

	radiusData := make([]Radius, 0)

	for _, v := range radiusServersQuery.List {
		radiusData = append(radiusData, Radius{Name: v.Name, ID: v.ID, IP: v.Primary.IP, Port: v.Primary.Port, Secret: v.Primary.SharedSecret})

	}
	return radiusData
}

func GetSiteName(ip string, domainsData []Domains) (string, string) {

	var radiusProfile string
	var zoneName string

	for _, domain := range domainsData {

		for _, zone := range domain.GetZones() {
			for _, radius := range zone.GetRadius() {
				if radius.IP == ip {
					radiusProfile = radius.Name
					zoneName = zone.Name
				}
			}
		}
	}

	return radiusProfile, zoneName

}

func GetAllRadiusServers(controller map[string]string) {

	type TextSearch struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	}

	type QueryFilter struct {
		TextSearch TextSearch `json:"fullTextSearch"`
		Limit      int64      `json:"limit"`
	}

	type RadiusServers struct {
		List []struct {
			ID         string `json:"id"`
			DomainID   string `json:"domainId"`
			Name       string `json:"name"`
			RadiusIP   string `json:"radiusIP"`
			RadiusPort int64  `json:"radiusPort"`
			ZoneUUID   string `json:"zoneUUID"`
		} `json:"list"`
	}

	textSearch := TextSearch{Type: "AND", Value: ""}
	filter := QueryFilter{Limit: 1000, TextSearch: textSearch}

	queryFilter, err := json.Marshal(filter)

	if err != nil {
		fmt.Println("error constructing query filter: ", err)
	}

	queryRadiusServersUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/query/services/aaaServer/auth"

	httpResp, _ := BuildHttpRequest(queryRadiusServersUrl, "POST", nil, bytes.NewBuffer(queryFilter), controller["accesstoken"], true)

	var radiusServers RadiusServers

	err = json.Unmarshal([]byte(httpResp.([]uint8)), &radiusServers)
	if err != nil {
		fmt.Println(err)
	}

	//return radiusServers

}
