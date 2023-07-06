package main

import (
	"log"
	"sync"
	"time"
)

func ParseEachController(cliArguments *Args) {
	var wg sync.WaitGroup

	for _, c := range cliArguments.Controllers {
		wg.Add(1)
		controller := map[string]string{"hostname": c.Hostname, "clustername": c.Clustername, "port": c.Port, "username": c.Username, "password": c.Password}

		go PollController(cliArguments, controller, &wg)

	}
	wg.Wait()
}

func PollController(cliArguments *Args, controller map[string]string, wg *sync.WaitGroup) {
	defer wg.Done()

	// loop indefinitely for each controller to setup continous polling

	for {
		Login(controller)

		if controller["accesstoken"] != "none" {

			domainsData := GetDomainsData(controller)

			onlineTotal, offlineTotal := GetApTotals(controller)
			log.Printf("Polled AP counts on %s", controller["hostname"])

			domainClientCounts := GetClientCountsDomain(controller, domainsData)
			totalClientCounts := GetClientCount(controller, "")
			ClientCountsExporter(cliArguments, controller["hostname"], totalClientCounts, domainClientCounts)
			log.Printf("Polled client counts on %s", controller["hostname"])

			totalLicenseCount, usedLicenseCount := GetLicenseTotals(controller)
			log.Printf("Polled license counts on %s", controller["hostname"])

			totalLicenseMetric, usedLicenseMetric := Exporter(cliArguments, controller["hostname"], "licenseCount")
			onlineApMetric, offlineApMetric := Exporter(cliArguments, controller["hostname"], "apCount")

			controllerNodeStates, overallClusterState := GetClusterState(controller)

			ClusterStateExporter(cliArguments, controller["hostname"], controllerNodeStates[controller["clustername"]], overallClusterState)

			totalLicenseMetric.Set(float64(totalLicenseCount))
			usedLicenseMetric.Set(float64(usedLicenseCount))
			onlineApMetric.Set(float64(onlineTotal))
			offlineApMetric.Set(float64(offlineTotal))

			log.Printf("exported metrics for  %s", controller["hostname"])

			alarmTotal := GetActiveAlarms(controller, domainsData)

			AlarmStateExporter(cliArguments, controller["hostname"], alarmTotal)

		}

		log.Printf("Polling %s in %v seconds", controller["hostname"], cliArguments.PollingPeriod)
		time.Sleep(time.Duration(cliArguments.PollingPeriod) * time.Second)
	}
}
