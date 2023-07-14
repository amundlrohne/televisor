package models

type resultType string

const (
	Vector resultType = "vector"
)

type PrometheusAPIResponse struct {
	Status string           `json:"status"`
	Data   PrometheusResult `json:"data"`
}

type PrometheusResult struct {
	ResultType resultType              `json:"resultType"`
	Result     []PromtheusVectorResult `json:"result"`
}

type PromtheusVectorResult struct {
	Metric PrometheusVectorResultMetric `json:"metric"`
	Value  []interface{}                `json:"value"` // [int, string]
}

type PrometheusVectorResultMetric struct {
	Name string `json:"name"`
}

type PrometheusContainerMetric map[string]string
