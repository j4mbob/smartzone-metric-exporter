package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
)

type Args struct {
	MetricHost    string `json:"metrichost"`
	MetricPort    string `json:"metricport"`
	PollingPeriod int    `json:"pollingperiod"`
	Controllers   []struct {
		Hostname    string `json:"hostname"`
		Clustername string `json:"clustername"`
		Port        string `json:"port"`
		Username    string `json:"username"`
		Password    string `json:"password"`
	} `json:"controllers"`
}

func main() {

	arguments := new(Args)
	ArgParse(arguments)
	ParseEachController(arguments)

}

func ArgParse(arguments *Args) {

	metricHostPtr := flag.String("metrichost", "http://grafana.networks-util.ask4.net", "remote metric store host to use")
	metricPortPtr := flag.String("metricport", "8428", "remote metric store port to use")
	pollingPeriodPtr := flag.Int("pollingperiod", 30, "polling period in seconds")

	configPtr := flag.String("loadconfig", "none", "load json config file. defaults to poller-config.json")

	flag.Parse()

	arguments.MetricHost = *metricHostPtr
	arguments.MetricPort = *metricPortPtr
	arguments.PollingPeriod = *pollingPeriodPtr

	if *configPtr != "none" {
		log.Printf("loading JSON config: %s\n", *configPtr)
		LoadConfig(*configPtr, arguments)
	}

}

func LoadConfig(configFile string, arguments *Args) {

	jsonFile, err := os.Open(configFile)

	if err != nil {
		log.Println(err)
		os.Exit(1)

	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal([]byte(byteValue), &arguments)

}
