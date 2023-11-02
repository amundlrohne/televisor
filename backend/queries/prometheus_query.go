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

func prometheusQuery(addr string, query string) []byte {
	res, err := http.Get(addr + "/api/v1/query?query=" + url.PathEscape(query))
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

func PrometheusContainerCPUQuantile(addr string) models.PrometheusContainerMetric {
	const quantile = `sum(quantile_over_time(0.997, rate(container_cpu_usage_seconds_total[30s])[1h:]) * 100) by (name)`

	response := prometheusQuery(addr, quantile)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerCPUStdev(addr string) models.PrometheusContainerMetric {
	const std = `sum(stddev_over_time(rate(container_cpu_usage_seconds_total[30s])[1h:])) by (name)`

	response := prometheusQuery(addr, std)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerCPUMean(addr string) models.PrometheusContainerMetric {
	const mean = `sum(avg_over_time(rate(container_cpu_usage_seconds_total[30s])[1h:]) * 100) by (name)`

	response := prometheusQuery(addr, mean)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryQuantile(addr string) models.PrometheusContainerMetric {
	const quantile = `sum by (name) ((quantile_over_time(0.997, container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes) * 100)`

	response := prometheusQuery(addr, quantile)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryStdev(addr string) models.PrometheusContainerMetric {
	const stdev = `sum by (name) (stddev_over_time(container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes)`

	response := prometheusQuery(addr, stdev)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerMemoryMean(addr string) models.PrometheusContainerMetric {
	const mean = `sum by (name) ((avg_over_time(container_memory_usage_bytes[1h]) / on() group_left() machine_memory_bytes) * 100)`

	response := prometheusQuery(addr, mean)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerNetworkOutput(addr string) models.PrometheusContainerMetric {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_output_total[15s]) / 1024)`
	response := prometheusQuery(addr, query)

	var promResult models.PrometheusAPIResponse
	json.Unmarshal(response, &promResult)

	return convertApiResponseToMap(promResult)
}

func PrometheusContainerNetworkInput(addr string) models.PrometheusContainerMetric {
	const query = `sum by(name) (podman_container_info{name!~".+infra"} * on(id) group_right(name) rate(podman_container_net_input_total[15s]) / 1024)`
	response := prometheusQuery(addr, query)

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
