package models

type AnnotationType string

const (
	Criticality AnnotationType = "CRITICALITY"
	Dependence  AnnotationType = "DEPENDENCE"
	Megaservice AnnotationType = "MEGASERVICE"
	Greedy      AnnotationType = "GREEDY"
	Cyclic      AnnotationType = "CYCLIC"
)

type YChartLevel string

const (
	OperationLevel   YChartLevel = "OPERATION"
	ServiceLevel     YChartLevel = "SERVICE"
	ApplicationLevel YChartLevel = "APPLICATION"
)

type Annotation struct {
	Services       []string
	Operations     []string
	Message        string
	AnnotationType AnnotationType
	YChartLevel    YChartLevel
}
