package transport

import (
	"regexp"
	"store/pkg/store/infrastructure/prometheus"
	"strings"
)

const PathPrefix = "/api/v1/"

func NewEndpointLabelCollector() prometheus.EndpointLabelCollector {
	return endpointLabelCollector{}
}

type endpointLabelCollector struct {
}

func (e endpointLabelCollector) EndpointLabelForURI(uri string) string {
	if strings.HasPrefix(uri, PathPrefix) {
		r, _ := regexp.Compile("^" + PathPrefix + "user/[a-f0-9-]+$")
		if r.MatchString(uri) {
			return specUserEndpoint
		}
	}
	return uri
}
