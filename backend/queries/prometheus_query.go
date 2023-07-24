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

func PrometheusContainerCPUQuantile() models.PrometheusContainerMetric {
	const quantile = `sum(quantile_over_time(0.997, rate(container_cpu_usage_seconds_total[30s])[1h:]) * 100) by (name)`

	response := prometheusQuery(quantile)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerCPUStdev() models.PrometheusContainerMetric {
	const std = `sum(stddev_over_time(rate(container_cpu_usage_seconds_total[30s])[1h:])) by (name)`

	response := prometheusQuery(std)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerCPUMean() models.PrometheusContainerMetric {
	const mean = `sum(avg_over_time(rate(container_cpu_usage_seconds_total[30s])[1h:]) * 100) by (name)`

	response := prometheusQuery(mean)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryQuantile() models.PrometheusContainerMetric {
	const quantile = `sum by (name) ((quantile_over_time(0.997, container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes) * 100)`

	response := prometheusQuery(quantile)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryStdev() models.PrometheusContainerMetric {
	const stdev = `sum by (name) (stddev_over_time(container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes)`

	response := prometheusQuery(stdev)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryMean() models.PrometheusContainerMetric {
	const mean = `sum by (name) ((avg_over_time(container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes) * 100)`

	response := prometheusQuery(mean)

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
