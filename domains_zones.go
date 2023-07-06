package main

import (
	"encoding/json"
	"fmt"
)

type Domains struct {
	ID             string
	ApCount        int64
	Name           string
	ParentDomainID string
	SubDomainCount int64
	ZoneCount      int64
	Zones          []Zone
}

type Radius struct {
	Name   string
	ID     string
	IP     string
	Port   int64
	Secret string
}

type Zone struct {
	ID            string
	Name          string
	RadiusServers []Radius
}

type ZoneValues interface {
	GetZones() []Zone
}

type RadiusValues interface {
	GetRadius() []Radius
}

func (d Domains) GetZones() []Zone {
	return d.Zones
}

func (z Zone) GetRadius() []Radius {
	return z.RadiusServers
}

type DomainsQuery struct {
	List []struct {
		ID             string `json:"id"`
		ApCount        int64  `json:"apCount"`
		CreateDatetime string `json:"createDatetime"`
		CreatedBy      string `json:"createdBy"`
		DomainType     string `json:"domainType"`
		Name           string `json:"name"`
		ParentDomainID string `json:"parentDomainId"`
		SubDomainCount int64  `json:"subDomainCount"`
		ZoneCount      int64  `json:"zoneCount"`
	}
}

func GetDomainsData(controller map[string]string) []Domains {

	queryDomainsUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/domains"

	params := make(map[string]string)
	params["listSize"] = "1000"

	httpResp, _ := BuildHttpRequest(queryDomainsUrl, "GET", params, nil, controller["accesstoken"], true)

	var domainsQuery DomainsQuery

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &domainsQuery)
	if err != nil {
		fmt.Println(err)
	}

	domainsData := GetZonesData(controller, domainsQuery)

	return domainsData

}

func GetZonesData(controller map[string]string, domainsQuery DomainsQuery) []Domains {

	type ZonesQuery struct {
		List []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		}
	}

	var zonesQuery ZonesQuery

	domainsData := make([]Domains, 0)

	for _, v := range domainsQuery.List {

		params := make(map[string]string)
		params["listSize"] = "1000"

		queryZonesUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/rkszones"

		params["domainId"] = v.ID

		httpResp, _ := BuildHttpRequest(queryZonesUrl, "GET", params, nil, controller["accesstoken"], true)

		err := json.Unmarshal([]byte(httpResp.([]uint8)), &zonesQuery)
		if err != nil {
			fmt.Println(err)
		}

		zoneData := make([]Zone, 0)

		for _, v := range zonesQuery.List {
			radiusData := GetZoneRadiusServer(controller, v.ID)
			zoneData = append(zoneData, Zone{v.ID, v.Name, radiusData})

		}
		domainsData = append(domainsData, Domains{ID: v.ID, ApCount: v.ApCount, Name: v.Name, ParentDomainID: v.ParentDomainID, SubDomainCount: v.SubDomainCount, ZoneCount: v.ZoneCount, Zones: zoneData})

	}

	return domainsData

}
