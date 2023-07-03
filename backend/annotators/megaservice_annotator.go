package annotators

import (
	"log"

	"github.com/amundlrohne/televisor/models"
	"github.com/jaegertracing/jaeger/model"
)

func MegaserviceAnnotator(spans []model.Span) models.Annotation {
	span := spans[0]

	log.Println("This is a span: %v", span)

	return models.Annotation{}
}
