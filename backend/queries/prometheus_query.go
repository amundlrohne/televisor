package queries

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/amundlrohne/televisor/models"
)

func prometheusQuery(query string) []byte {
	res, err := http.Get("http://localhost:9090/api/v1/query?query=" + url.PathEscape(query))
	if err != nil {
		fmt.Printf("error making http request: %s\n", err)
		os.Exit(1)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error reading response body: %s\n", err)
		os.Exit(1)
	}

	return resBody
}

func PrometheusContainerCPU() models.PrometheusAPIResponse {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_cpu_seconds_total[15s]))`

	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	fmt.Printf("client: response body: %+v\n", promResult)
	return promResult
}

func PrometheusContainerMemory() models.PrometheusAPIResponse {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) podman_container_mem_usage_bytes / 1024 /100)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	fmt.Printf("client: response body: %+v\n", promResult)
	return promResult
}

func PrometheusContainerNetworkOutput() models.PrometheusAPIResponse {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_output_total[15s]) / 1024)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	fmt.Printf("client: response body: %+v\n", promResult)
	return promResult
}

func PrometheusContainerNetworkInput() models.PrometheusAPIResponse {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_input_total[15s]]) / 1024)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	fmt.Printf("client: response body: %+v\n", promResult)
	return promResult
}
