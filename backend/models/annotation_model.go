package models

type AnnotationType string

const (
	Criticality            AnnotationType = "CRITICALITY"
	Dependence             AnnotationType = "DEPENDENCE"
	Megaservice            AnnotationType = "MEGASERVICE"
	Greedy                 AnnotationType = "GREEDY"
	InappropriateIntimacy  AnnotationType = "INAPPROPRIATE_INTIMACY"
	Cyclic                 AnnotationType = "CYCLIC"
	OverUtilizationCPU     AnnotationType = "OVER_UTILIZATION_CPU"
	UnderUtilizationCPU    AnnotationType = "UNDER_UTILIZATION_CPU"
	OverUtilizationMemory  AnnotationType = "OVER_UTILIZATION_MEMORY"
	UnderUtilizationMemory AnnotationType = "UNDER_UTILIZATION_MEMORY"
)

type YChartLevel string

const (
	OperationLevel   YChartLevel = "OPERATION"
	ServiceLevel     YChartLevel = "SERVICE"
	ApplicationLevel YChartLevel = "APPLICATION"
)

type AnnotationLevel string

const (
	Critical AnnotationLevel = "CRITICAL"
	Warning  AnnotationLevel = "WARNING"
	Info     AnnotationLevel = "INFO"
)

type Annotation struct {
	Services            []string        `json:"services"`
	Operations          []string        `json:"operations"`
	InitiatingOperation string          `json:"initiatingOperation"`
	Message             string          `json:"message"`
	AnnotationType      AnnotationType  `json:"annotationType"`
	YChartLevel         YChartLevel     `json:"yChartLevel"`
	AnnotationLevel     AnnotationLevel `json:"annotationLevel"`
	Recomendation       Recommendation  `json:"recommendation"`
}
