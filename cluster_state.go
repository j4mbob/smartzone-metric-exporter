package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type ClusterState struct {
	ClusterName                string `json:"clusterName"`
	ClusterRole                string `json:"clusterRole"`
	ClusterState               string `json:"clusterState"`
	CurrentNodeID              string `json:"currentNodeId"`
	CurrentNodeName            string `json:"currentNodeName"`
	ManagementServiceStateList []struct {
		ManagementServiceState string `json:"managementServiceState"`
		NodeID                 string `json:"nodeId"`
		NodeName               string `json:"nodeName"`
	} `json:"managementServiceStateList"`
	NodeStateList []struct {
		NodeID    string `json:"nodeId"`
		NodeName  string `json:"nodeName"`
		NodeState string `json:"nodeState"`
	} `json:"nodeStateList"`
}

func GetClusterState(controller map[string]string) (map[string]map[string]map[string]string, string) {

	queryApUrl := "https://" + controller["hostname"] + ":" + controller["port"] + "/wsg/api/public/v10_0/cluster/state"

	var cluster ClusterState

	httpResp, _ := BuildHttpRequest(queryApUrl, "GET", nil, nil, controller["accesstoken"], true)

	err := json.Unmarshal([]byte(httpResp.([]uint8)), &cluster)
	if err != nil {
		fmt.Println("error: ", err)
	}

	controllerNodeStates, overallClusterState := GetNodeStates(cluster)

	return controllerNodeStates, overallClusterState

}

func GetNodeStates(cluster ClusterState) (map[string]map[string]map[string]string, string) {

	controllerNodeStates := make(map[string]map[string]map[string]string)

	controllerNodeStates[cluster.ClusterName] = make(map[string]map[string]string)

	for _, node := range cluster.NodeStateList {
		controllerNodeStates[cluster.ClusterName][node.NodeName] = make(map[string]string)
		controllerNodeStates[cluster.ClusterName][node.NodeName]["nodeState"] = node.NodeState
	}

	for _, node := range cluster.ManagementServiceStateList {
		controllerNodeStates[cluster.ClusterName][node.NodeName]["managementServiceState"] = node.ManagementServiceState
	}

	log.Printf("Polled cluster state on %s", cluster.ClusterName)

	return controllerNodeStates, cluster.ClusterState

}
