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

func error() {
	fmt.Fprintf(os.Stderr, "Usage: get_heapster_metrics.go <object (pod/node)> <namespace> <object name> <metric name>\n")
	os.Exit(1)
}

func main() {
	if len(os.Args) < 2 {
		error()
	}
	object := os.Args[1]
	if (object == "pod" || object == "cluster") && (len(os.Args) < 4) {
		error()
	}
	if (object != "metrics") && (len(os.Args) < 3) {
		error()
	}

	baseURL := "http://10.11.4.19:32652/api/v1/model/"
	url := ""
	if object == "pod" {
		ns := os.Args[2]
		name := os.Args[3]
		url = baseURL + "namespaces/" + ns + "/pods/" + name + "/metrics"

	} else if object == "cluster" {
		ns := os.Args[2]
		name := os.Args[3]
		url = baseURL + "namespaces/" + ns + "/nodes/" + name + "/metrics"

	} else if object == "metrics" {
		url = baseURL + "metrics"
	} else {
		ns := os.Args[2]
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
