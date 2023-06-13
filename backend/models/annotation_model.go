package models

type AnnotationType string

const (
	Criticality AnnotationType = "CRITICALITY"
	Dependency  AnnotationType = "DEPENDENCY"
)

type Annotation struct {
	Services       []string
	AnnotationType AnnotationType
}
