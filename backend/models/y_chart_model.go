package models

type YChartModel struct {
	Annotations []Annotation                `json:"annotations"`
	Services    map[string]TelevisorService `json:"services"`
	Operations  Operations                  `json:"operations"`
}
