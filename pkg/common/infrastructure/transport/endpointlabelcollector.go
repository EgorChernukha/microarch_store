package transport

import (
	"store/pkg/common/infrastructure/prometheus"
)

func NewEndpointLabelCollector() prometheus.EndpointLabelCollector {
	return endpointLabelCollector{}
}

type endpointLabelCollector struct {
}

func (e endpointLabelCollector) EndpointLabelForURI(uri string) string {
	return uri
}
