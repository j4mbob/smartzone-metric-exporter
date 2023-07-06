package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func Login(controller map[string]string) {

	type Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Timezone string `json:"timeZoneUtcOffset"`
	}

	creds := Credentials{controller["username"], controller["password"], "+00:00"}
	bodyParams, err := json.Marshal(creds)
	if err != nil {
		log.Println(err)
		os.Exit(1)

	}

	requestTokenUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/session"

	log.Printf("attempting to login to %s", controller["hostname"])

	_, cookieData := BuildHttpRequest(requestTokenUrl, "POST", nil, bytes.NewBuffer(bodyParams), "", false)

	var accessToken string

	for _, c := range cookieData {
		if c.Name == "JSESSIONID" {
			accessToken = c.Value

		}
	}

	if accessToken != "" {
		log.Println("obtained access token")
		controller["accesstoken"] = accessToken
	} else {
		log.Println("could not obtain access token")
		controller["accesstoken"] = "none"
	}

}

func BuildHttpRequest(requestUrl string, method string, urlParams map[string]string, bodyParams *bytes.Buffer, accessToken string, returnJson bool) (interface{}, []*http.Cookie) {

	var request *http.Request
	var err error

	var params url.Values

	if urlParams != nil {
		params = url.Values{}
		for param, value := range urlParams {
			params.Add(param, value)
		}
		requestUrl = requestUrl + "?" + params.Encode()
	}

	if bodyParams == nil {
		request, err = http.NewRequest(method, requestUrl, nil)
	} else {
		request, err = http.NewRequest(method, requestUrl, bodyParams)
	}

	if err != nil {
		log.Printf("An Error Occured %v", err)
		return false, nil
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if accessToken != "" {
		bearer := fmt.Sprintf("JSESSIONID=%s", accessToken)
		request.Header.Set("Cookie", bearer)
		request.Header = http.Header{
			"Content-Type": {"application/json; charset=UTF-8"},
			"Cookie":       {bearer},
		}

	}

	data, cookieData := ExecuteHttpRequest(request, returnJson)

	return data, cookieData

}

func ExecuteHttpRequest(request *http.Request, returnJson bool) (interface{}, []*http.Cookie) {

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("error: %s", err)
		return false, nil
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error: %s", err)
		return false, nil
	}

	cookieData := response.Cookies()

	if returnJson {
		if response.StatusCode != 200 {
			log.Printf("http error: %v  ", response.StatusCode)
			return false, nil
		}
		return body, cookieData
	}

	bodyData := JsonToMapData(body)

	if response.StatusCode != 200 {
		log.Printf("http error: %v details: %v", response.StatusCode, bodyData["error"])
		return false, nil
	}

	return bodyData, cookieData
}

func JsonToMapData(jsondata []byte) map[string]interface{} {

	var data map[string]interface{}

	err := json.Unmarshal([]byte(jsondata), &data)
	if err != nil {
		log.Printf("error: %s", err)
	}

	return data
}
