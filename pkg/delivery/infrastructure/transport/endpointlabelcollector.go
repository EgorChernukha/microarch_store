package transport

import (
	"regexp"
	"strings"

	"store/pkg/common/infrastructure/prometheus"
)

const PathPrefix = "/api/v1/"

func NewEndpointLabelCollector() prometheus.EndpointLabelCollector {
	return endpointLabelCollector{}
}

type endpointLabelCollector struct {
}

func (e endpointLabelCollector) EndpointLabelForURI(uri string) string {
	if strings.HasPrefix(uri, PathPrefix) {
		r, _ := regexp.Compile("^" + PathPrefix + "delivery/[a-f0-9-]+$")
		if r.MatchString(uri) {
			return specDeliveryEndpoint
		}
	}
	return uri
}
