package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type PodIntervalData struct {
	Average    int64 `json:"average,omitempty"`
	Percentile int64 `json:"percentile,omitempty"`
	Max        int64 `json:"max,omitempty"`
}

type PodResources struct {
	Minute PodIntervalData `json:"minute,omitempty"`
	Hour   PodIntervalData `json:"hour,omitempty"`
	Day    PodIntervalData `json:"day,omitempty"`
}

type PodStats struct {
	CpuLimit      PodResources `json:"cpu-limit,omitempty"`
	CpuUsage      PodResources `json:"cpu-usage,omitempty"`
	MemoryLimit   PodResources `json:"memory-limit,omitempty"`
	MemoryUsage   PodResources `json:"memory-usage,omitempty"`
	MemoryWorking PodResources `json:"memory-working,omitempty"`
}

type PodFullStats struct {
	Uptime int64    `json:"uptime,omitempty"`
	Stats  PodStats `json:"stats,omitempty"`
}

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: %s <object (pod/node)> <namespace> <object name>\n", os.Args[0])
		os.Exit(1)
	}
	object := os.Args[1]
	ns := os.Args[2]
	name := os.Args[3]

	baseURL := "http://10.11.4.19:32652/api/v1/model/"
	url := ""
	if object == "pod" {
		url = baseURL + "namespaces/" + ns + "/pods/" + name + "/metrics"

	} else if object == "cluster" {
		url = baseURL + "namespaces/" + ns + "/nodes/" + name + "/metrics"

	} else if object == "list" {
		url = baseURL + "metrics"
	} else {
		url = baseURL + "namespaces/" + ns + "/pods"
	}
	if len(os.Args) > 4 {
		metric := os.Args[4]
		url = url + "/" + metric
	}

	response, err := http.Get(url)
	fmt.Printf("Requested url %s\n", url)

	var FullStats PodFullStats
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &FullStats)
	// prettyParsedBody, err := json.MarshalIndent(FullStats, "", "  ")
	if err != nil {
		fmt.Printf("ErrorY :%v\n", err)
	} else {
		// fmt.Printf("Fullstats:", FullStats)
		fmt.Printf("Full body:\n")
		fmt.Printf(string(body))
	}

	_, err = io.Copy(os.Stdout, response.Body)
	if err != nil {
		log.Fatal(err)
	}

}
