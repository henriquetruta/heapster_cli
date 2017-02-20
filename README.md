# Basic heapster CLI client

Simple heapster API parser to CLI commands, made in Golang.

This all assumes:

* you have Go installed and configured
* heapster API is accessible through NodePort in port 32652

## Using

This scripts converts some API calls to simple CLI usage. The usage is described as follows:

To run it, just do:

`go run get_heapster_metrics.go <ip> <object> <namespace> <object name> <metric>`

Where:

* `<ip>` is the IP to access the InfluxDB. If it's a NodePort service, then it's the host's IP;
* `<object>` the kind of object you want to query on. Supports `pod` and `node` values. If `metrics` is provided, all metrics are listed. If any other value is providded, all pods inside the given namespace are listed;
* `<namespace>` the namespace the object belongs to;
* `<object name>` the name of the pod of node to query on;
* `<metric name>` (optional) the metric to be queried on. If not provided, all metrics available to the object are shown;

### Valid executions

* `go run get_heapster_metrics.go list default`: Lists all pods in namespace `default`
* `go run get_heapster_metrics.go metrics`: Lists all available metrics;
* `go run get_heapster_metrics.go pod default mypod`: Lists all available metrics for pod called `mypod` in namespace `default`;
* `go run get_heapster_metrics.go pod default mypod cpu/usage`: Returns cpu usage for pod called `mypod` in namespace `default`;