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

func PrometheusContainerCPU() models.PrometheusContainerMetric {
	const query = `sum by (name) (rate(container_cpu_usage_seconds_total[30s]) * 100)`
	// const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_cpu_seconds_total[15s]))`

	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemory() models.PrometheusContainerMetric {
	const query = `sum by (name) ((container_memory_usage_bytes / on() group_left() machine_memory_bytes) * 100)`
	// const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) podman_container_mem_usage_bytes / 1024 /100)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerNetworkOutput() models.PrometheusContainerMetric {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_output_total[15s]) / 1024)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerNetworkInput() models.PrometheusContainerMetric {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_input_total[15s]) / 1024)`
	response := prometheusQuery(query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func convertApiResponseToMap(response models.PrometheusAPIResponse) models.PrometheusContainerMetric {
	containerMetricMap := make(map[string]string)

	for _, m := range response.Data.Result {
		containerMetricMap[m.Metric.Name] = m.Value[len(m.Value)-1].(string)
	}

	return containerMetricMap
}
