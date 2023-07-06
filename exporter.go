package main

import (
	"log"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

func Exporter(arguments *Args, hostname string, metric string) (*metrics.FloatCounter, *metrics.FloatCounter) {
	s := metrics.NewSet()

	if metric == "licenseCount" {
		totalLicenseCount := s.NewFloatCounter("totalLicenseCount")
		usedLicenseCount := s.NewFloatCounter("usedLicenseCount")

		err := s.InitPush(arguments.MetricHost+":"+arguments.MetricPort+"/api/v1/import/prometheus/metrics/job/licenseCount/instance/"+hostname, time.Duration(arguments.PollingPeriod)*time.Second, "")
		if err != nil {
			log.Fatalln("error: ", err)
		}

		return totalLicenseCount, usedLicenseCount

	} else if metric == "apCount" {
		onlineTotal := s.NewFloatCounter("onlineApTotal")
		offlineTotal := s.NewFloatCounter("offlineApTotal")

		err := s.InitPush(arguments.MetricHost+":"+arguments.MetricPort+"/api/v1/import/prometheus/metrics/job/apCount/instance/"+hostname, time.Duration(arguments.PollingPeriod)*time.Second, "")
		if err != nil {
			log.Fatalln("error: ", err)
		}

		return onlineTotal, offlineTotal
	}

	return nil, nil
}

func ClusterStateExporter(arguments *Args, hostname string, controllerNodeStates map[string]map[string]string, overallClusterState string) {
	s := metrics.NewSet()

	err := s.InitPush(arguments.MetricHost+":"+arguments.MetricPort+"/api/v1/import/prometheus/metrics/job/clusterState/instance/"+hostname, time.Duration(arguments.PollingPeriod)*time.Second, "")
	if err != nil {
		log.Fatalln("error: ", err)
	}

	for k, v := range controllerNodeStates {
		managementState := s.NewCounter(`managementServiceState{node="` + k + `"}`)
		nodeState := s.NewCounter(`nodeState{node="` + k + `"}`)

		if v["managementServiceState"] == "In_Service" {
			managementState.Set(1)
		} else {
			managementState.Set(2)

		}

		if v["nodeState"] == "In_Service" {
			nodeState.Set(1)
		} else {
			nodeState.Set(2)

		}
	}

	clusterState := s.NewCounter("overallClusterState")

	if overallClusterState == "In_Service" {
		clusterState.Set(1)
	} else {
		clusterState.Set(2)
	}
}

func AlarmStateExporter(arguments *Args, hostname string, activeAlarms map[string]int) {

	s := metrics.NewSet()

	err := s.InitPush(arguments.MetricHost+":"+arguments.MetricPort+"/api/v1/import/prometheus/metrics/job/activeAlarmCount/instance/"+hostname, time.Duration(arguments.PollingPeriod)*time.Second, "")
	if err != nil {
		log.Fatalln("error: ", err)
	}

	for k, v := range activeAlarms {
		alarmRadiusServer := s.NewCounter(`radiusAlarm{server="` + k + `"}`)
		alarmRadiusServer.Set(uint64(v))

	}

}

func ClientCountsExporter(arguments *Args, hostname string, clientCounts int, domainClientCounts map[string]int) {
	s := metrics.NewSet()

	err := s.InitPush(arguments.MetricHost+":"+arguments.MetricPort+"/api/v1/import/prometheus/metrics/job/clientCounts/instance/"+hostname, time.Duration(arguments.PollingPeriod)*time.Second, "")
	if err != nil {
		log.Fatalln("error: ", err)
	}

	clientCount := s.NewFloatCounter("totalClientCount")
	clientCount.Set(float64(clientCounts))

	for k, v := range domainClientCounts {
		domainClientCount := s.NewCounter(`clientCounts{domain="` + k + `"}`)
		domainClientCount.Set(uint64(v))
	}

}
