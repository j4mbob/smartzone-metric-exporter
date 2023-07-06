package main

import (
	"encoding/json"
	"fmt"
)

type Licenses struct {
	List []struct {
		CapacityControlLicenseCount struct {
			TotalCount int64 `json:"totalCount"`
			UsedCount  int64 `json:"usedCount"`
		} `json:"capacityControlLicenseCount"`
		LicenseTypeDescription string `json:"licenseTypeDescription"`
	} `json:"list"`
}

func GetLicenseTotals(controller map[string]string) (int64, int64) {

	var totalLicenseCount, usedLicenseCount int64

	licenseData := GetLicenseData(controller)

	for _, licenseType := range licenseData.List {
		if licenseType.LicenseTypeDescription == "AP Capacity License" {
			totalLicenseCount, usedLicenseCount = licenseType.CapacityControlLicenseCount.TotalCount, licenseType.CapacityControlLicenseCount.UsedCount
		}
	}

	return totalLicenseCount, usedLicenseCount
}

func GetLicenseData(controller map[string]string) Licenses {

	licenseUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/licensesSummary"
	data, _ := BuildHttpRequest(licenseUrl, "GET", nil, nil, controller["accesstoken"], true)

	var licenseData Licenses

	err := json.Unmarshal([]byte(data.([]uint8)), &licenseData)
	if err != nil {
		fmt.Println("error: ", err)
	}

	return licenseData

}
